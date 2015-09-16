package gopool
import (
	"testing"
	_ "log"
	"reflect"
)

func sqr(x int) int {
	return x * x
}

func inSlice(a interface{}, list interface{}) bool {
	for i:=0; i<reflect.ValueOf(list).Len(); i++ {
		if reflect.DeepEqual(reflect.ValueOf(a).Interface(),reflect.ValueOf(list).Index(i).Interface()) {
			return true
		}
	}
	return false
}

func avg(x []int) float64 {
	var total float64
	for _, i := range x {
		total += float64(i)
	}
	return total / float64(len(x))
}
func dumb(x int) {

}

func TestSqr(t *testing.T) {
	slice := []int{3,4,5,6,7,8}
	resultSlice := []int64{9,16,25,36,49,64}
	pool := New(10)
	pool.Map(sqr, slice)

	result := pool.Join()
	rresult := reflect.ValueOf(result)

	if reflect.TypeOf(result).Kind() != reflect.Slice {
		t.Errorf("Map doesn't return slice! It returns %s", reflect.TypeOf(result).Kind())
	} else {
		for i:=0; i<rresult.Len(); i++ {
			if !inSlice(rresult.Index(i).Int(), resultSlice){
				t.Errorf("%d is not equal of any %v", rresult.Index(i).Int(), resultSlice)
			}
		}
	}
}

func TestAvg(t *testing.T) {
	slice := [][]int{{3, 4, 5, 6}, {1, 2, 3, 4, 5}}
	resultSlice := []float64{4.5, 3}
	pool := New(1)
	pool.Map(avg, slice)
	result := pool.Join()
	rresult := reflect.ValueOf(result)

	if reflect.TypeOf(result).Kind() != reflect.Slice {
		t.Errorf("Map doesn't return slice! It returns %s", reflect.TypeOf(result).Kind())
	} else {
		for i:=0; i<rresult.Len(); i++ {
			if !inSlice(rresult.Index(i).Float(), resultSlice){
				t.Errorf("%f is not equal of any %v", rresult.Index(i).Float(), resultSlice)
			}
		}
	}
}
func TestDumb(t *testing.T) {
	slice := []int{3,4,5,6,7,8}
	pool := New(10)
	pool.Map(dumb, slice)
	pool.Join()
}