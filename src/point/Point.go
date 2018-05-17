package point

type Point struct {

	X float64
	Y float64
	Ip int16

}

func NewPoint(x,y float64,ip int16) Point {
	return Point{x,y,ip}
}
