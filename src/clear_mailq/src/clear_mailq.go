/*
	This Code Is Overly Commented For Training and Auditing Purposes
*/

package main

// Let's Use These Modules
import(
	"fmt"
	"os"
	"path/filepath"
)

func main() {

	// Check If The Directory Even Exists
	_, err := os.Stat("/var/spool/clientmqueue")
	if os.IsNotExist(err) {
		fmt.Println("directory doesn't exist?: /var/spool/clientmqueue")
		os.Exit(1)
	}

	// It Does!  Oh Good.  Let's Change To It
	os.Chdir("/var/spool/clientmqueue")


	// Globbing Is Fast, So If There Are A Lot Of Files To Delete
	// It Won't Take An Extensively Long Time
	matches, err := filepath.Glob("*")

	// Always Check For Errors	
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// This Gets Us The Path... After The Chdir It Should Always Be "/var/spool/clientmqueue"
	// But This Is Important To Make Sure We Can Audit What Was Deleted, And Make Sure We Are
	// Only Deleting The Things We Want To
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))


	// Make Certain We Are Only Deleting Things From "/var/spool/clientmqueue"
	if dir != "/var/spool/clientmqueue" {
		fmt.Println("Dir Is NOT /var/spool/clientmqueue, cowardly backing away")
		os.Exit(1)
	}


	// If There Are More Than 100 Files, Don't Bother Listing All Of them
	printflag := true
	if len(matches) > 100 {
		printflag = false
	}


	// If There Are No Files, Then Just Exit
	if len(matches) <= 0 {
		fmt.Println("Nothing to delete")
		return
	}

	// Loop Through All Of The Found Files That We Globbed Earlier, And Remove Them.
	// If It's Less Than 100, Display Each File Removed.  Otherwise, Don't Do That.
	for i := 0; i < len(matches); i++ {
		os.Remove(matches[i])
		if printflag {
			fmt.Println("Removing: " + dir + "/" + matches[i])
		}
	}


	// Change The Way We Output The Final Result If There Are More Than 100 Files
	if printflag {
		fmt.Println("successfully deleted mail")
	} else {
		fmt.Println("successfully deleted over 100 emails")
	}


	// Cleanly Exit.  Exit Code 0 Is Typically Success.  Exit Code 1 Is An Error.
	os.Exit(0)
}

