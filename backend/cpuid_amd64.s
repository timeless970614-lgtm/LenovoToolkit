// +build amd64

#include "textflag.h"

// func cpuid(leaf uint32, eax, ebx, ecx, edx *uint32)
TEXT ·cpuid(SB), NOSPLIT, $0-40
	MOVL leaf+0(FP), AX
	CPUID
	MOVQ eax+8(FP), DI
	MOVL AX, 0(DI)
	MOVQ ebx+16(FP), DI
	MOVL BX, 0(DI)
	MOVQ ecx+24(FP), DI
	MOVL CX, 0(DI)
	MOVQ edx+32(FP), DI
	MOVL DX, 0(DI)
	RET
