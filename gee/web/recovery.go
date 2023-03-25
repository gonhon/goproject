/*
 * @Author: gaoh
 * @Date: 2023-03-25 18:21:26
 * @LastEditTime: 2023-03-25 18:51:11
 */
package web

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

func Recovery() HandlerFunc {
	return func(ctx *Context) {
		//异常捕获
		defer func() {
			if err := recover(); err != nil {
				msg := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(msg))
				ctx.Fail(http.StatusInternalServerError, "Internal Server Error")
			}

		}()
		ctx.Next()
	}
}

func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
