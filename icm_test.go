package icm

import (
	"fmt"
	"image/png"
	"os"
	"testing"
	"time"
)

// load image
func TestICM(t *testing.T) {
	orig_name := "noisyImage.png"
	file, err := os.Open(orig_name)
	if err != nil {
		os.Exit(1)
	}
	defer file.Close() // close file when exiting..
	img, err := png.Decode(file)
	if err != nil {
		t.Error(err)
	}

	// use these H, Beta, Eta values
	h, b, e := 0.0, 1.0, 1.0

	// instantiate Unidirectional Graphical Model
	u := NewUGM(img, h, b, e)

	begin := time.Now()

	fmt.Printf("Initial E = %.0f\n", u.E())

	not_noisy := "i=%02d_h=%.1f_b=%.1f_e=%.1f_E_%.0f.png"
	for i := 1; i < 20; i++ {
		u.ICM(1)
		fmt.Printf("i = %2d: E = %.0f\n", i, u.E())
		u.WriteToFile(fmt.Sprintf(not_noisy, i, h, b, e, u.E()), img)
	}

	total := time.Now().Sub(begin).Nanoseconds()
	// total := end - begin
	// fmt.Printf("%d - %d = %d", end, begin, total)

	fmt.Printf("-----------------------\nTime: %4d Milliseconds\n", total/1e6)
}
