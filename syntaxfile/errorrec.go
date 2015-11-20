/*
 * Copyright(C) 2015 Simon Schmidt
 * 
 * This Source Code Form is subject to the terms of the
 * Mozilla Public License, v. 2.0. If a copy of the MPL
 * was not distributed with this file, You can obtain one at
 * http://mozilla.org/MPL/2.0/.
 * 
 */

package syntaxfile

type SyntaxError struct{
	Position interface{}
	What string
}

type ErrorRecorder struct{
	Errors []SyntaxError
}
func (e *ErrorRecorder)Backup() (b ErrorRecorder){
	if e==nil { return }
	b = *e
	return
}
func (e *ErrorRecorder)Restore(b ErrorRecorder){
	if e==nil { return }
	*e = b
}
func (e *ErrorRecorder)Add(err SyntaxError){
	if e==nil { return }
	e.Errors = append(e.Errors,err)
}
func (e *ErrorRecorder)AddErr(p interface{},w string){
	if e==nil { return }
	e.Errors = append(e.Errors,SyntaxError{p,w})
}

