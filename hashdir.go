package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strings"
	"fmt"
	"time"
)

// Make generate hash of all files and they paths for specified directory.
func Make(dir string) (string, error) {
	var fileHash_sha256 string
	var fileHash_md5 string
	buf := make([]byte, 1024*1024)

	bigErr := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
			
				//combined sha256 and md5 file hash
				fileHash_sha256 = ""
				fileHash_md5 = ""
				
					// Handle hashing big files.
		// Source: https://stackoverflow.com/q/60328216/1722542

				f, err := os.Open(path)
				if err != nil {
					//return err
					if os.IsPermission(err) {
						fmt.Printf("%s,AccessDenied-InvalidPermissions,,,\n", path)
					} else {
						fmt.Printf("%s,UnknownFileIssue,,,\n", path)
					}
				} else {


					defer func() {
						_ = f.Close()
					}()

					
					sha256h, err := selectHash("sha256")
					if err != nil {
						return err
					}
					
					md5h, err := selectHash("md5")
					if err != nil {
						return err
					}

					for {
						bytesRead, err := f.Read(buf)
						if err != nil {
							if err != io.EOF {
								return err
							}
							_, err = sha256h.Write(buf[:bytesRead])
							if err != nil {
								return err
							}
							_, err = md5h.Write(buf[:bytesRead])
							if err != nil {
								return err
							}
							break
						}
						_, err = sha256h.Write(buf[:bytesRead])
						if err != nil {
							return err
						}
						_, err = md5h.Write(buf[:bytesRead])
						if err != nil {
							return err
						}
					}
					

					fileHash_sha256 = hex.EncodeToString(sha256h.Sum(nil))
					fileHash_md5 = hex.EncodeToString(md5h.Sum(nil))
					
					fmt.Printf("%s,%s,%s,%d,%s\n", path, fileHash_sha256, fileHash_md5, info.Size(), info.ModTime().Format("2006-01-02 15:04:05 UTC"))
				}

			}
			return nil
		})
			//
			
			
			
//previous code used the generic has function, but this is slower as it reads the file twice			
/*

				fileHash_sha256, err := hashFile(path, "sha256")
				if err != nil {
					//return err
					if os.IsPermission(err) {
						fmt.Printf("%s,AccessDenied-InvalidPermissions,,,\n", path)
                         		} else {
                         			fmt.Printf("%s,UnknownFileIssue,,,\n", path)
                         		}

				} else {
				
					//continue with next hash if it didn't error
				
					fileHash_md5, err := hashFile(path, "md5")
					if err != nil {
						//shouldn't happen because already checked above, but comment out to ensure unexpected stop
						//return err
					}

					//log.Println(path, fileHash, info.ModTime(), info.Size())
					fmt.Printf("%s,%s,%s,%d,%s\n", path, fileHash_sha256, fileHash_md5, info.Size(), info.ModTime().Format("2006-01-02 15:04:05 UTC"))
					
					//fileHash = fileHash_sha256
					//endHash = endHash + pathHash + fileHash
				}
*/				
				
			
		

	return "", bigErr
}

func selectHash(hashType string) (hash.Hash, error) {
	switch strings.ToLower(hashType) {
	case "md5":
		return md5.New(), nil
	case "sha1":
		return sha1.New(), nil
	case "sha256":
		return sha256.New(), nil
	case "sha512":
		return sha512.New(), nil
	}
	return nil, errors.New("Unknown hash: " + hashType)
}

func hashData(data string, hashType string) (string, error) {
	h, err := selectHash(hashType)
	if err != nil {
		return "", err
	}

	_, err = h.Write([]byte(data))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}



//general hashfile
func hashFile(path string, hashType string) (string, error) {
	// Handle hashing big files.
	// Source: https://stackoverflow.com/q/60328216/1722542

	f, err := os.Open(path)
	if err != nil {
		//this will error if there is no file access
		return "", err
	}


	defer func() {
		_ = f.Close()
	}()

	buf := make([]byte, 1024*1024)
	h, err := selectHash(hashType)
	if err != nil {
		return "", err
	}
	

	for {
		bytesRead, err := f.Read(buf)
		if err != nil {
			if err != io.EOF {
				return "", err
			}
			_, err = h.Write(buf[:bytesRead])
			if err != nil {
				return "", err
			}
			break
		}
		_, err = h.Write(buf[:bytesRead])
		if err != nil {
			return "", err
		}
	}
	

	fileHash := hex.EncodeToString(h.Sum(nil))
	
	return fileHash, nil
}


// Main function
func main() {
  
    var dir string      
    dir, err := os.Getwd() // For read access.
    if err != nil {
	os.Exit(1)
    }
    
    if len(os.Args) > 1 {
      dir = os.Args[1]
    }
  
    // Finding the time
    fmt.Println("Time: ", time.Now().Unix())
    fmt.Println("Dir: ", dir)
    
    fmt.Println("path, fileHash_sha256, fileHash_md5, filesize, modified_time")    
    
    Make(dir)
    
}
