package parse

import (
	"testing"
)

func TestAccept(t *testing.T) {
	// Accept(0) will always match :)
	p0 := Accept(0)
	if n, r := p0(""); !n.Matched || n.Content != "" || n.Nodes != nil || r != "" {
		t.Error("Invalid result for Accept empty match test", n)
	}
	// Test with a non-zero amount of bytes
	p := Accept(1)
	if n, r := p(""); n.Matched || n.Content != "" || n.Nodes != nil || r != "" {
		t.Error("Invalid result for Accept no-match test", n)
	}
	if n, r := p("a"); !n.Matched || n.Content != "a" || n.Nodes != nil || r != "" {
		t.Error("Invalid result for Accept match test", n)
	}
}

func TestK(t *testing.T) {
	p := Accept(1)
	k := K(p)
	if n, r := k("aaa"); !n.Matched || n.Content != "aaa" || n.Nodes == nil || r != "" {
		t.Error("Invalid result for K* match test", n)
	}
	if n, r := k(""); !n.Matched || n.Content != "" || n.Nodes != nil || r != "" {
		t.Error("Invalid result for K* no-match test", n)
	}
}

func TestSeq(t *testing.T) {
	p1, p2 := Accept(1), Accept(1)
	s := Seq(p1, p2)
	if n, r := s("aa"); !n.Matched || n.Content != "aa" || n.Nodes == nil || r != "" {
		t.Error("Invalid result for Seq positive test", n)
	}
	if n, r := s("a"); n.Matched || n.Content != "a" || n.Nodes == nil || r != "" {
		t.Error("Invalid result for Seq partial match test", n)
	}
	if n, r := s(""); n.Matched || n.Content != "" || n.Nodes == nil || r != "" {
		t.Error("Invalid result for Seq no-match test", n)
	}
}

func TestAny(t *testing.T) {
	p, p0 := Accept(1), Accept(0)
	a := Any(p, p0)
	if n, r := a(""); !n.Matched || n.Content != "" || n.Nodes != nil || r != "" {
		t.Error("Invalid result for Any match test", n)
	}
	if n, r := a("aa"); !n.Matched || n.Content != "a" || n.Nodes != nil || r != "a" {
		t.Error("Invalid result for Any match test", n)
	}
}

func TestDefer(t *testing.T) {
	d, dp := Defer()
	*dp = Accept(1)
	if n, r := d(""); n.Matched || n.Content != "" || n.Nodes != nil || r != "" {
		t.Error("Invalid result for Defer no-match test", n)
	}
	if n, r := d("a"); !n.Matched || n.Content != "a" || n.Nodes != nil || r != "" {
		t.Error("Invalid result for Defer match test", n)
	}
}

func TestString(t *testing.T) {
	s := String("aa")
	if n, r := s("aa"); !n.Matched || n.Content != "aa" || n.Nodes != nil || r != "" {
		t.Error("Invalid result for String match test", n)
	}
	if n, r := s("ab"); n.Matched || n.Content != "" || n.Nodes == nil || r != "ab" {
		t.Error("Invalid result for String no-match test", n, r, len(r))
	}
}

func TestDigit(t *testing.T) {
	d := Digit()
	if n, r := d("a"); n.Matched || n.Content != "" || n.Nodes == nil || r != "a" {
		t.Error("Invalid result for Digit no-match test", n)
	}
	if n, r := d("1"); !n.Matched || n.Content != "1" || n.Nodes != nil || r != "" {
		t.Error("Invalid result for Digit match test", n)
	}
}
