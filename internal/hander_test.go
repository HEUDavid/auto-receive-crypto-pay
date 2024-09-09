package internal

import "testing"

func TestFSM(t *testing.T) {
	if err := ReceiptFSM.Draw("receipt.svg"); err != nil {
		t.Fatal(err)
	}
}
