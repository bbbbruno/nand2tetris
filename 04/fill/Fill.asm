// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel;
// the screen should remain fully black as long as the key is pressed. 
// When no key is pressed, the program clears the screen, i.e. writes
// "white" in every pixel;
// the screen should remain fully clear as long as no key is pressed.

// Put your code here.
(LOOP)
    @SCREEN
    D=A
    @R0
    M=D      // R0にSCREENの開始アドレスを入れる

    @KBD
    D=M
    @BLACK   // KEYが押されているかどうかで分岐
    D;JNE
    @WHITE
    D;JEQ
(BLACK)
    @R0
    D=M
    @KBD
    D=D-A
    @LOOP
    D;JEQ    // R0のアドレスがKBDのアドレスまで到達していたらループを抜ける

    @R0
    A=M      // R0のアドレス先に切り替える
    M=-1     // 切り替え先の値を全てオンにする

    @R0
    M=M+1    // 次の16ビットワードに移って繰り返し
    @BLACK
    0;JMP
(WHITE)
    @R0
    D=M
    @KBD
    D=D-A
    @LOOP
    D;JEQ    // R0のアドレスがKBDのアドレスまで到達していたらループを抜ける

    @R0
    A=M      // R0のアドレス先に切り替える
    M=0      // 切り替え先の値を全てオフにする

    @R0
    M=M+1
    @WHITE   // 次の16ビットワードに移る
    0;JMP
