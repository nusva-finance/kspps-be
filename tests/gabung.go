package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// 1. Tentukan nama file hasil gabungan
	outputFilename := "gabungan_kspps_be.txt"

	// 2. Tentukan folder mana saja yang ingin digabung (sesuaikan dengan struktur folder Anda)
	targetFolders := []string{"../cmd", "../config", "../handlers", "../middleware", "../models", "../repositories", "../routes", "../services"}

	// 3. Tentukan ekstensi file yang ingin diambil
	allowedExtensions := []string{".go", ".sql", ".json"}

	outFile, err := os.Create(outputFilename)
	if err != nil {
		fmt.Println("Gagal membuat file output:", err)
		return
	}
	defer outFile.Close()

	for _, folder := range targetFolders {
		// Pengecekan jika folder ada
		if _, err := os.Stat(folder); os.IsNotExist(err) {
			fmt.Printf("Melewati %s (folder tidak ditemukan)\n", folder)
			continue
		}

		err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			
			// Jika ini bukan folder, cek ekstensinya
			if !info.IsDir() {
				ext := filepath.Ext(path)
				for _, allowedExt := range allowedExtensions {
					if ext == allowedExt {
						// Tulis header penanda lokasi file
						header := fmt.Sprintf("\n\n============================================================\nFILE: %s\n============================================================\n\n", path)
						outFile.WriteString(header)

						// Buka dan copy isi file-nya
						content, err := os.ReadFile(path)
						if err != nil {
							outFile.WriteString(fmt.Sprintf("// Error membaca file: %v\n", err))
						} else {
							outFile.Write(content)
						}
						break
					}
				}
			}
			return nil
		})

		if err != nil {
			fmt.Printf("Error saat menelusuri folder %s: %v\n", folder, err)
		}
	}
	fmt.Printf("Selesai! File berhasil digabung dan disimpan di: %s\n", outputFilename)
}