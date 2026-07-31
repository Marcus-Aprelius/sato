package main

import (
	"fmt"
	"log"

	"sato/internal/testdb"
)

func main() {
	if err := testdb.WriteDemoDatabase("secrets.kdbx", testdb.DefaultPassword); err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅ Test database created: secrets.kdbx")
	fmt.Println("   Password: sato")
	fmt.Println("   Entries:")
	fmt.Println("   - DB_PASSWORD")
	fmt.Println("   - API_KEY")
	fmt.Println("   - NGINX_PASSWORD")
	fmt.Println("   - test/test1")
	fmt.Println("   - test/test2")
	fmt.Println("   Empty groups:")
	fmt.Println("   - empty-group")
	fmt.Println("   - test/empty-nested")
}
