/*
   应该Wrap这个error，抛给上层。
   因为一般顶层调用者，在拿到error的时候，属于业务逻辑关键功能，需要定位并处理，
   需要拿到更多关于该error的详细信息（包错误信息，调用栈等）
   解决办法:
   import "github.com/pkg/errors" 的Wrap方法对错误进行上下文包装，
   并携带原始错误信息，尽量保留完整的调用栈,方便调用者定位问题
 */
package main
import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
)

func GetSql() error{
	return errors.Wrap(sql.ErrNoRows,"GetSql failed")
}

func Call() error{
	return errors.WithMessage(GetSql(),"Call failed")
	//return errors.WithStack(GetSql())
}

func main()  {
	err := Call()
	if err != nil{
		if errors.Cause(err) == sql.ErrNoRows {
			fmt.Printf("data not found: %v\n",err)
			fmt.Printf("%+v\n",err)
			return
		}
	}
}
