package mutexscope

import (
	"go/types"

	"github.com/gostaticanalysis/analysisutil"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"
)

const doc = "mutexscope mutexscope finds sync.Mutex which likely forgotten .Unlock()."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "mutexscope",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		buildssa.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	syncMutex := analysisutil.TypeOf(pass, "sync", "Mutex")
	if syncMutex == nil {
		return nil, nil
	}
	lockFuncObj := analysisutil.MethodOf(syncMutex, "Lock")
	unlockFuncObj := analysisutil.MethodOf(syncMutex, "Unlock")
	if lockFuncObj == nil || unlockFuncObj == nil {
		return nil, nil
	}

	funcs := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA).SrcFuncs

	for _, f := range funcs {
		for _, b := range f.Blocks {
			i, val := findLock(lockFuncObj, b)
			if i == -1 {
				continue
			}

			seekedBlock := make(map[*ssa.BasicBlock]struct{})
			if canReachReturnWithoutUnlock(unlockFuncObj, b, i+1, val, seekedBlock) {
				pass.Reportf(b.Instrs[i].Pos(), "Should Unlock inside function.")
			}
		}
	}

	return nil, nil
}

func findLock(lockFunc *types.Func, b *ssa.BasicBlock) (int, ssa.Value) {
	for index, instr := range b.Instrs {
		callInstr, _ := instr.(*ssa.Call)
		if callInstr == nil {
			continue
		}

		callee := callInstr.Call.StaticCallee()
		if callee == nil {
			continue
		}

		funcObj := callee.Object()

		if len(callInstr.Call.Args) < 1 {
			continue
		}
		arg := callInstr.Call.Args[0]

		if funcObj == lockFunc && arg != nil {
			return index, arg
		}
	}
	return -1, nil
}

func canReachReturnWithoutUnlock(unlockFunc *types.Func, b *ssa.BasicBlock, start int, v ssa.Value, seeked map[*ssa.BasicBlock]struct{}) bool {
	if _, ok := seeked[b]; ok {
		return false
	}
	seeked[b] = struct{}{}

	instrs := b.Instrs

	var returnInstr *ssa.Return
	for _, instr := range instrs {
		r, _ := instr.(*ssa.Return)
		if r != nil {
			returnInstr = r
		}
	}

	for i := start; i < len(instrs); i++ {
		instr := instrs[i]
		if isUnlock(unlockFunc, v, instr) {
			return false
		}
		if dfsCanReachReturnWithoutUnlock(unlockFunc, instr, v, seeked) {
			return true
		}
		if isFuncCallWhichIncludesUnlock(unlockFunc, instr, v, seeked) {
			return false
		}
	}
	return returnInstr != nil
}

func dfsCanReachReturnWithoutUnlock(unlockFunc *types.Func, instr ssa.Instruction, v ssa.Value, seeked map[*ssa.BasicBlock]struct{}) bool {
	succs := instr.Block().Succs
	for _, succ := range succs {
		if canReachReturnWithoutUnlock(unlockFunc, succ, 0, v, seeked) {
			return true
		}
	}
	return false
}

func isFuncCallWhichIncludesUnlock(unlockFunc *types.Func, instr ssa.Instruction, v ssa.Value, seeked map[*ssa.BasicBlock]struct{}) bool {
	callInstr, _ := instr.(*ssa.Call)
	if callInstr == nil {
		return false
	}

	callee := callInstr.Call.StaticCallee()
	if callee != nil {
		blocks := callee.Blocks
		// nilのときパッケージ外
		if blocks != nil {
			if !canReachReturnWithoutUnlock(unlockFunc, blocks[0], 0, v, seeked) {
				return true
			}
		}
	}
	return false
}

func isUnlock(unlockFunc *types.Func, v ssa.Value, instr ssa.Instruction) bool {
	callInstr, _ := instr.(*ssa.Call)
	if callInstr == nil {
		deferInstr, _ := instr.(*ssa.Defer)
		if deferInstr == nil {
			return false
		}
		return isUnlockCallCommon(unlockFunc, v, deferInstr.Call)
	}
	return isUnlockCallCommon(unlockFunc, v, callInstr.Call)
}

func isUnlockCallCommon(unlockFunc *types.Func, v ssa.Value, cc ssa.CallCommon) bool {
	callee := cc.StaticCallee()
	if callee == nil {
		return false
	}

	funcObj := callee.Object()
	if len(cc.Args) < 1 {
		return false
	}
	arg := cc.Args[0]

	if funcObj == unlockFunc && arg == v {
		return true
	}
	return false
}
