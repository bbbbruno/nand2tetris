package tokenizer

import (
	"bufio"
	"io"
)

type Tokenizer interface {
	HasMoreTokens() bool
	Advance()
	CurrentToken() Token
}

type tokenizer struct {
	*bufio.Reader
	currentToken Token
	nextToken    Token
	buf          []byte
}

func New(r io.Reader) Tokenizer {
	b := bufio.NewReader(r)
	return &tokenizer{b, nil, nil, make([]byte, 0)}
}

// 次のトークンが取得できるかどうか確認する。
func (tkz *tokenizer) HasMoreTokens() bool {
	t, err := tkz.parse()
	if t == nil || err != nil {
		return false
	}

	tkz.nextToken = t
	return true
}

// 次のトークンを現在のトークンにセットする。
func (tkz *tokenizer) Advance() {
	tkz.currentToken = tkz.nextToken
}

// 入力を解析して次のトークンを生成する。
func (tkz *tokenizer) parse() (Token, error) {
	// bufが空のとき。
	if len(tkz.buf) == 0 {
		// 次の一文字を読み込み、その次の一文字も先読みする。
		b, err := tkz.ReadByte()
		if err != nil {
			return nil, err
		}
		bs, err := tkz.Peek(1)
		if err != nil {
			return nil, err
		}
		nextb := bs[0]

		// 次の文字が空白文字ならスキップ。
		if isSpace(b) {
			return tkz.parse()
		}
		// 次の文字からコメントならコメント部分を丸ごとスキップ。
		if isComment(b, nextb) {
			b, _ := tkz.ReadByte()
			if err := tkz.skipComments(b); err != nil {
				return nil, err
			}
			return tkz.parse()
		}
		// 次の文字から文字列が開始するならSTRING_CONSTのトークンを生成。
		if isStringConst(b) {
			s, err := tkz.readString()
			if err != nil {
				return nil, err
			}

			return &token{STRING_CONST, s}, nil
		}
		// 次の文字が終端記号ならSYMBOLのトークンを生成。
		if isSymbol(b) {
			return &token{SYMBOL, string(b)}, nil
		}

		// どれにも一致しないなら読み込んだ文字をbufに溜めて次に行く。
		tkz.buf = append(tkz.buf, b)
		return tkz.parse()
	} else { // bufにすでに文字が溜まっているとき。
		// 次の一文字を先読みする。
		bs, err := tkz.Peek(1)
		if err != nil {
			return nil, err
		}
		nextb := bs[0]

		// 先読みした文字がIdentifierとして使える文字（半角英数字とアンダースコア）ならbufに溜めて次へ行く。
		if isIdentifier(string(nextb)) {
			b, _ := tkz.ReadByte()
			tkz.buf = append(tkz.buf, b)
			return tkz.parse()
		} else { // Identifierとして使えない文字ならbufをクリアしてトークンを生成する。
			buf := tkz.buf
			tkz.clearBuf()
			if isKeyword(buf) { // bufの文字列がキーワードとして認識可能ならKEYWORDのトークンを生成する。
				return &token{KEYWORD, string(buf)}, nil
			} else if isIntConst(buf) { // bufの文字列が整数値として変換できるならINT_CONSTのトークンを生成する。
				return &token{INT_CONST, string(buf)}, nil
			} else { // それ以外はIDENTIFIERのトークンを生成する。
				return &token{IDENTIFIER, string(buf)}, nil
			}
		}
	}
}

func (tkz *tokenizer) CurrentToken() Token {
	return tkz.currentToken
}

// bufを空にする。
func (tkz *tokenizer) clearBuf() {
	tkz.buf = []byte{}
}

// コメント部分を丸ごとスキップする。
func (tkz *tokenizer) skipComments(b byte) error {
	switch b {
	case '*':
		bs, err := tkz.Peek(1)
		if err != nil {
			return err
		}
		nextb := bs[0]
		delim := byte('*')

		if nextb == '*' {
			_, _ = tkz.ReadByte()
			for b != '/' {
				if _, err := tkz.ReadBytes(delim); err != nil {
					return err
				}
				if b, err = tkz.ReadByte(); err != nil {
					return err
				}
			}
		} else {
			if _, err := tkz.ReadBytes(delim); err != nil {
				return err
			}
			if _, err := tkz.ReadByte(); err != nil {
				return err
			}
		}
	case '/':
		if _, _, err := tkz.ReadLine(); err != nil {
			return err
		}
	}

	return nil
}

// 文字列を最後まで読み込む。
func (tkz *tokenizer) readString() (string, error) {
	str, err := tkz.ReadString('"')
	if err != nil {
		return "", err
	}

	return str[:len(str)-1], nil
}
