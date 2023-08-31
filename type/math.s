#include "textflag.h"

// func add(a, b int) int
TEXT ·add(SB), NOSPLIT, $0-24
    MOVQ a+0(FP), AX // 参数a
    MOVQ b+8(FP), BX // 参数b
    ADDQ BX, AX
    MOVQ AX, ret+16(FP)
    RET

// func sub(a, b int) int
TEXT ·sub(SB), NOSPLIT, $0-24
    MOVQ a+0(FP), AX // 参数a
    MOVQ b+8(FP), BX // 参数b
    SUBQ BX, AX
    MOVQ AX, ret+16(FP)
    RET

// func mul(a, b int) int
TEXT ·mul(SB), NOSPLIT, $0-24
    MOVQ a+0(FP), AX // 参数a
    MOVQ b+8(FP), BX // 参数b
    IMULQ BX, AX
    MOVQ AX, ret+16(FP)
    RET

