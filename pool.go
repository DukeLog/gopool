package gopool

import (
	"log"
	"reflect"
	"sync"
)


type Pool struct {
	size        int
	wg          *sync.WaitGroup
	iChannel    chan interface{}
	oChannel    chan interface{}
	return_type reflect.Type
}

func (this *Pool) worker(f interface{}, i int) {
	defer this.wg.Done()
	vf := reflect.ValueOf(f)

	for i := range this.iChannel {
		vi := reflect.ValueOf(i)
		if this.return_type != nil {
			ret := vf.Call([]reflect.Value{vi})[0]
			this.oChannel <- ret.Interface()
		} else {
			vf.Call([]reflect.Value{vi})
			this.oChannel <- nil
		}
	}
}

func New(size int) *Pool {
	p := new(Pool)
	p.size = size
	p.wg = new(sync.WaitGroup)
	p.iChannel = make(chan interface{})
	return p
}

func (this *Pool) Map(f interface{}, i interface{}) {
	fval := reflect.ValueOf(f)
	ftype := fval.Type()
	if ftype.Kind() != reflect.Func {
		log.Panicf("`f` must be type %s , but get %s", reflect.Func, ftype.Kind())
	}
	if ftype.NumIn() != 1 {
		log.Panicf("`f` should have only one parameter, but get %d parameters", ftype.NumIn())
	}
	if ftype.NumOut() > 1 {
		log.Panicf("`f` should return not more than one value but it returns %d values", ftype.NumOut())
	}

	for i := 0; i < this.size; i++ {
		this.wg.Add(1)
		go this.worker(f, i)
	}
	value_i := reflect.ValueOf(i)
	if ftype.NumOut() == 1 {
		this.return_type = ftype.Out(0)
	} else {
		this.return_type = nil
	}

	this.oChannel = make(chan interface{}, value_i.Len())
	for ii := 0; ii < value_i.Len(); ii++ {
		this.iChannel <- value_i.Index(ii).Interface()
	}
}

func (this *Pool) Join() interface{} {
	close(this.iChannel)
	this.wg.Wait()
	close(this.oChannel)
	if this.return_type != nil {
		dynamic_slice := reflect.SliceOf(this.return_type)
		ret := reflect.MakeSlice(dynamic_slice, 0, 0)
		for i := range this.oChannel {
			ret = reflect.Append(ret, reflect.ValueOf(i))
		}
		return ret.Interface()
	} else {
		return nil
	}

}