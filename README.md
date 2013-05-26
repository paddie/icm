icm
===

Image de-noising using ICM covered in Bishop 8.3.3 using energy function from equation 8.4.2

### Example:
```Go
func main() {
  orig_name := "noisyImage.png"
	file, err := os.Open(orig_name)
	if err != nil {
		os.Exit(1)
	}
	defer file.Close() // close file when exiting..
	img, err := png.Decode(file)
	if err != nil {
		panic(err)
	}
	h, b, e := 0.0, 1.0, 1.0

	u := NewUGM(img, h, b, e)

	fmt.Printf("DIM: (%d,%d)\n", len(u.image), len(u.image[0]))
	iterations := 1
	now := time.Now().Nanosecond()
	fmt.Printf("Before %d iterations: E: %.0f\n", iterations, u.E())
	u.ICM(iterations)
	now = time.Now().Nanosecond() - now
	fmt.Printf("After  %d iterations: E: %.0f\nTime: %d Milliseconds\n", iterations, u.E(), now/1e6)

	otherimg := image.NewGray16(img.Bounds())
	for x := 0; x < u.X(); x++ {
		for y := 0; y < u.Y(); y++ {
			if u.Node(x, y).X() == -1 {
				otherimg.SetGray16(x, y, color.Black)
			} else {
				otherimg.SetGray16(x, y, color.White)
			}
		}
	}
	not_noisy := "ICM__i=%d_h=%.1f_b=%.1f_e=%.1f_E_%.0f.png"
	newImg, err := os.Create(fmt.Sprintf(not_noisy, iterations, h, b, e, u.E()))
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
	defer newImg.Close() // close when exiting..

	/* can only enconde in PNG - 'jpeg' package is still too.. young ;)*/
	if err := png.Encode(newImg, otherimg); err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
}
```
### Output
```
DIM: (300,300)
Before 1 iterations: E: -195150
After  1 iterations: E: -238352
Time: 25 Milliseconds
```

### TODO:
- Instead of iterations argument, make a stopping criteria 
