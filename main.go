package main

// El modulo flag nos ayuda a definir los flags que aceptara nuestro programa y a parsearlos de forma automatica desde la linea de comandos,
// por ejemplo, si ejecutamos nuestro programa con la opcion -p, el valor de -p sera almacenado en la variable pattern
import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"

	"golang.org/x/exp/constraints"
)

func main() {
	// Sintaxis: flag.String(nombre, valor por defecto, mensaje de ayuda)
	// * Filter flags
	flagPattern := flag.String("p", "", "filter by pattern")
	flagAll := flag.Bool("a", false, "show all files including hidden files")
	flagNumberRecords := flag.Int("n", 0, "number of records to show")
	// * Order flags
	hasOrderByTime := flag.Bool("t", false, "sort by time, oldest first")
	hasOrderBySize := flag.Bool("s", false, "sort by size, smallest first")
	hasOrderReverse := flag.Bool("r", false, "reverse order")

	// Con .Parse() parseamos los flags que se pasaron por linea de comandos, parsear significa que se extrae el valor de los flags y se almacena en la variable correspondiente
	flag.Parse()
	// fmt.Println("pattern:", *flagPattern)
	// fmt.Println("all:", *flagAll)
	// fmt.Println("number of records:", *flagNumberRecords)
	// fmt.Println("order by time:", *hasOrderByTime)
	// fmt.Println("order by size:", *hasOrderBySize)
	// fmt.Println("reverse order:", *hasOrderReverse)

	// Implemantacion de la logica de los argumentos de los diferentes flags

	// * flag.Args() retorna un slice de strings con los argumentos que no son flags, es decir, con los argumentos
	// flag.Args()
	path := flag.Arg(0)

	// * si el usuario no especifica un path, se asume que el path es el directorio actual "."
	if path == "" {
		path = "."
	}

	dirs, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	fs := []file{}

	for _, dir := range dirs {
		isHidden := isHidden(dir.Name(), path)

		if isHidden && !*flagAll {
			continue
		}

		if *flagPattern != "" {
			// con "(?i)" le indicamos a la expresion regular que no sea case sensitive
			isMatched, err := regexp.MatchString("(?i)"+*flagPattern, dir.Name())
			if err != nil {
				panic(err)
			}
			if !isMatched {
				continue
			}
		}

		f, err := getFile(dir, isHidden)
		if err != nil {
			panic(err)
		}

		fs = append(fs, f)
	}

	// Ordenamiento
	if !*hasOrderByTime || !*hasOrderBySize {
		orderByName(fs, *hasOrderReverse)
	}

	if !*hasOrderByTime || *hasOrderBySize {
		orderBySize(fs, *hasOrderReverse)
	}

	if *hasOrderByTime {
		orderByTime(fs, *hasOrderReverse)
	}

	if *flagNumberRecords == 0 || *flagNumberRecords > len(fs) {
		*flagNumberRecords = len(fs)
	}
	printList(fs, *flagNumberRecords)
}

func mySort[T constraints.Ordered](i, j T, isReverse bool) bool {
	if isReverse {
		return i > j
	}
	return i < j
}

func orderByTime(files []file, isReverse bool) {
	sort.SliceStable(files, func(i, j int) bool {
		return mySort(
			files[i].modificationTime.Unix(),
			files[j].modificationTime.Unix(),
			isReverse,
		)
	})
}

func orderByName(files []file, isReverse bool) {
	sort.SliceStable(files, func(i, j int) bool {
		// if isReverse {
		// 	return strings.ToLower(files[i].name) > strings.ToLower(files[j].name)
		// }
		// return strings.ToLower(files[i].name) < strings.ToLower(files[j].name)
		return mySort(
			strings.ToLower(files[i].name), 
			strings.ToLower(files[j].name), 
			isReverse,
		)
	})
}

func orderBySize(files []file, isReverse bool) {
	sort.SliceStable(files, func(i, j int) bool {
		// if isReverse {
		// 	return files[i].size > files[j].size
		// }
		// return files[i].size < files[j].size
		return mySort(
			files[i].size,
			files[j].size,
			isReverse,
		)
	})
}

func printList(fs []file, nRecords int) {
	for _, file := range fs[:nRecords] {
		style := mapStyleByFileType[file.fileType]
		fmt.Printf("%s %s %s %10d %s %s %s %s\n", file.mode, file.userName, file.groupName, file.size, file.modificationTime.Format(time.DateTime), style.icon, file.name, style.symbol)
	}
}

func getFile(dir fs.DirEntry, isHidden bool) (file, error) {
	info, err := dir.Info()
	if err != nil {
		return file{}, fmt.Errorf("dir.Info(): %v", err)
	}

	f:= file {
		name: dir.Name(),
		fileType: 0,
		isDir: dir.IsDir(),
		isHidden: isHidden,
		userName: "cristian",
		groupName: "ATS",
		size: info.Size(),
		modificationTime: info.ModTime(),
		mode: info.Mode().String(),
	}
	setFile(&f)

	return f, nil
}

func setFile (f *file) file {
	switch {
	case isLink(*f):
		f.fileType = fileLink
	case f.isDir:
		f.fileType = fileDirectory
	case isExec(*f):
		f.fileType = fileExecutable
	case isCompress(*f):
		f.fileType = fileCompress
	case isImage(*f):
		f.fileType = fileImage
	default:
		f.fileType = fileRegular
	}
	return *f
}

func isLink(f file) bool {
	return strings.HasPrefix(strings.ToUpper(f.mode), "L")
}

func isExec(f file) bool {
	if runtime.GOOS == Windows {
		return strings.HasSuffix(f.name, exe)
	}
	return strings.Contains(f.mode, "x")
}

func isCompress(f file) bool {
	return strings.HasSuffix(f.name, zip) || 
				strings.HasSuffix(f.name, gz) || 
				strings.HasSuffix(f.name, tar) || 
				strings.HasSuffix(f.name, rar) || 
				strings.HasSuffix(f.name, deb)
}

func isImage(f file) bool {
	return strings.HasSuffix(f.name, png) || 
				strings.HasSuffix(f.name, jpg) || 
				strings.HasSuffix(f.name, gif)
}

func isHidden(fileName, basePath string) bool {
	return strings.HasPrefix(fileName, ".")
}