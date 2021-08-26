package main

import (
	"io"
	"os"
	"fmt"
	"time"

	"io/ioutil"
	"net/http"
	"path/filepath"
	"archive/tar"
	"compress/gzip"
)

func Log(message string) {
        current_time := time.Now().Local()
        t := current_time.Format("Jan 02 2006 03:04:05")

	fmt.Println(t + " - " + message)

        file, err := os.OpenFile("/tmp/auto_oracle.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil {
                return 
        }
        defer file.Close()

        _, err = file.WriteString(t + " - " + message + "\n")

        if err != nil {
                return 
        }

        return 
}

func Download(filepath string, url string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil	
}

func Untar(dst string, r io.Reader) error {

	tr := tar.NewReader(r)

	for {
		header, err := tr.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}

		target := filepath.Join(dst, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
			
			f.Close()
		}
	}
}

func UntarGZ(dst string, r io.Reader) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}

		target := filepath.Join(dst, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
			
			f.Close()
		}
	}
}

func Copy(src string, dst string) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dst, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing Database Name")
		return
	}

	dbname := os.Args[1]
	Log("Starting Install For: " + dbname)
	Log("Downloading auto_oracle_nonrac.tar.gz From http://pu-objectstore-01.cac.com:9000/oracle/auto_oracle_nonrac.tar.gz")

	Download("/root/auto_oracle_nonrac.tar.gz", "http://pu-objectstore-01.cac.com:9000/oracle/auto_oracle_nonrac.tar.gz")

	Log("Untaring /root/auto_oracle_nonrac.tar.gz to /root/auto_oracle")
	f, err := os.Open("/root/auto_oracle_nonrac.tar.gz")
	if err != nil {
		Log("Error Opening /root/auto_oracle_nonrac.tar.gz" + err.Error())
		return
	}

	err = UntarGZ("/root/", f)
	if err != nil {
		Log("Error Untaring /root/auto_oracle_nonrac.tar.gz: " + err.Error())
		return
	}
	f.Close()

	Log("Cleaning Up Our Mess: Deleting /root/auto_oracle_nonrac.tar.gz")
	err = os.Remove("/root/auto_oracle_nonrac.tar.gz")
	if err != nil {
		Log("Error Removing /root/auto_oracle_nonrac.tar.gz")
		return
	}

	Log("Extracting Template...")
	template, err := os.Open("/root/auto_oracle/template.tar")
	if err != nil {
		Log("Error Opening /root/auto_oracle/template.tar" + err.Error())
		return
	}

	err = Untar("/proj/app/", template)
	if err != nil {
		Log("Error Untaring /root/auto_oracle/template.tar: " + err.Error())
		return
	}

		
}

/*


echo 'stopping puppet' >> /tmp/auto_oracle.log
### Stop Puppet ###
service puppet stop


echo 'sed /etc/hosts' >> /tmp/auto_oracle.log
### Setup /etc/hosts ###
hst=$(hostname); sed -i "s/127.0.0.1.*//*127.0.0.1 $hst.cac.com ${hst} localhost.localdomain localhost/gi" /etc/hosts

#if [[ ! $hst =~ 01 ]]; then
#       echo "Exiting on $hst" 
#       echo "Exiting on $hst" >> /tmp/auto_oracle.log 
#       exit
#fi


echo 'register with spacewalk' >> /tmp/auto_oracle.log
### Register this shiz with spacewalk for yum updates to work ###
rhnreg_ks --serverUrl=https://pu-mtspcwlk-01.cac.com/XMLRPC --sslCACert=/usr/share/rhn/RHN-ORG-TRUSTED-SSL-CERT --activationkey=1-CAC --force



echo 'import keys' >> /tmp/auto_oracle.log
### Get Those Keys Yo ###
rpm --import /etc/pki/rpm-gpg/*


echo 'fix chrony' >> /tmp/auto_oracle.log
### Fix Chrony ###
sed -i 's/0.rhel.pool.ntp.org/ntp.cac.com/' /etc/chrony.conf


*/
