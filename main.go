package main

import (
  "os"
  "fmt"
  "time"
  "image"
  "image/png"
  "image/draw"

  "github.com/nfnt/resize"
)

func main () {
  // Load screenshot
  ss_file, err := os.Open("testscreenshot.png")
  if err != nil {
    fmt.Println("Error opening test screenshot: ", err)
  }
  ss_img, err := png.Decode(ss_file)
  if err != nil {
    fmt.Println("Error decoding test screenshot: ", err)
  }

  // Load map images into map of maps
  maps_map := make(map[string]image.Image)
  map_files := []string{"nCrypt1.png", "nCrypt2.png", "nCrypt3.png", "nCrypt4.png", "nCrypt5.png", "nCrypt6.png"}
  
  for _, filename := range map_files {
    file, err := os.Open(fmt.Sprintf("maps/%s",filename))
    if err != nil {
      fmt.Println("Error opening file: ", err)
      continue
    }
    img, err := png.Decode(file)
    if err != nil {
      fmt.Println("Error decoding image file: ", err)
    }
    maps_map[filename] = img
  }

  fmt.Println(maps_map)

  // Parse the minimap out of the screenshot
  tlx := 1668
  tly := 829

  minimap := ss_img.(interface {
    SubImage (r image.Rectangle) image.Image
  }).SubImage(image.Rect(tlx, tly, tlx+220, tly+220))

  // Calculate the perceptual hash for the minimap image.
  minimapHash := perceptualHash(minimap)

  // Normalize the minimap and the full maps. Preferably, shrink full maps by factor of 2
  starttime := time.Now()
  for name, pic := range maps_map {
    nw := float64(pic.Bounds().Dx()) / 2.05
    nh := float64(pic.Bounds().Dy()) / 2.05

    new_pic := resize.Resize(uint(nw), uint(nh), pic, resize.Lanczos3)
    maps_map[name] = new_pic
  }
  elapsed := time.Since(starttime)
  fmt.Printf("Lanczos3 resizing 6 crypt maps took: %s \n", elapsed)

  // Match the minimap image to a location on one of six full maps
  
