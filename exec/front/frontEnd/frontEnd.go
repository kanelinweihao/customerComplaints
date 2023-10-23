package frontEnd

import (
	// "fmt"
	"go.lwh.com/linweihao/customerComplaints/utils/err"
	// "go.lwh.com/linweihao/customerComplaints/utils/time"
	"go.lwh.com/linweihao/customerComplaints/exec/front/htmlSet"
)

func ExecFrontEnd() {
	defer err.ThrowError()
	htmlSet.SetHtml()
	return
}
