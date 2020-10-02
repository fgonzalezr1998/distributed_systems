package main

import (
  "fmt"
  "math"
)

func is_complex(a, b, c float32) bool {
  /*
   * Returns if the result of the ecuation is a complex number or not
   */

  return math.Pow(float64(b), float64(2.0)) < float64(4.0 * a + c)
}

func real_res(a, b float64, sqrt_arg float64) (r1, r2 float32) {
  /*
   * Returns both real result
   */

  r1 = float32((-b + math.Sqrt(sqrt_arg)) / (2.0 * a))
  r2 = float32((-b - math.Sqrt(sqrt_arg)) / (2.0 * a))

  return r1, r2
}

func print_real_res(a, b float64, sqrt_arg float64) {
  /*
   * Print real result
   */

  var res1, res2 float32

  res1, res2 = real_res(a, b, sqrt_arg)

  fmt.Println(res1, ", ", res2)
}

func print_complex_res(a, b float64, sqrt_arg float64) {
  /*
   * Print both complex results
   */

  var real_part, imag_part float32

  real_part = float32(-b / (2.0 * a))
  imag_part = float32(math.Sqrt(sqrt_arg))

  fmt.Println(real_part, " + ", imag_part, "j", ", ", real_part, " - ", imag_part)
}

func print_result(a, b, c float64) {
  /*
   * Print total result
   */

  var sqrt_arg float64
  var is_complex bool = false

  sqrt_arg = math.Pow(b, 2.0) - 4.0 * a * c
  if (sqrt_arg < 0) {
    is_complex = true
    sqrt_arg = -sqrt_arg
  }

  if (is_complex) {
    print_complex_res(a, b, sqrt_arg)
  } else {
    print_real_res(a, b, sqrt_arg)
  }
}

// Trial Consts:

const A, B, C float64 = 1.0, 4.0, 0.0

func main() {
  fmt.Print("Result: ")
  print_result(A, B, C)
}
