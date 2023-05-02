package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Stringer interface {
	String() string
}

func getUrl(md_key int) string {
	return "https://www.bilibili.com/bangumi/media/md" + strconv.Itoa(md_key)
}

func strBGM(md_key int, delimiter string) string {
	title := *bgm_data_map[md_key]
	url := getUrl(md_key)
	return fmt.Sprintf("%d%s%s%s%s", md_key, delimiter, title, delimiter, url)
}

func strlnBGM(md_key int, delimiter string) string {
	return strBGM(md_key, delimiter) + "\n"
}

func extractPageInfoByAPI(md_key int) error {
	resp, err := http.Get("http://api.bilibili.com/pgc/review/user?media_id=" + strconv.Itoa(md_key))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var jsonResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&jsonResponse)
	if err != nil {
		return err
	}

	var title string
	if result, ok := jsonResponse["result"].(map[string]interface{}); ok {
		if media, ok := result["media"].(map[string]interface{}); ok {
			if v, ok := media["title"].(string); ok {
				title = v
			} else {
				return errors.New("")
			}
		} else {
			return errors.New("")
		}
	} else {
		return errors.New("")
	}

	title_copy := new(string)
	*title_copy = title
	bgm_mutex.Lock()
	bgm_data_map[md_key] = title_copy
	bgm_mutex.Unlock()

	return nil
}

func saveData(path string, delimiter string) {
	cwd, err := os.Getwd()
	check(err)
	f, err := os.Create(filepath.Join(cwd, path))
	check(err)
	defer f.Close()

	f.Sync()

	w := bufio.NewWriter(f)
	for i := begin_idx; i < end_idx; i++ {
		if bgm_data_map[i] != nil {
			_, err = w.WriteString(strlnBGM(i, delimiter))
			check(err)
		}
	}

	w.Flush()
}

type WorkRoutine struct {
	idx   int
	begin int
	end   int
	sep   string
}

func (r *WorkRoutine) startRoutine(wg *sync.WaitGroup, bar *progressbar.ProgressBar) {
	defer wg.Done()
	for key := r.begin; key < r.end; key++ {
		bar.Add(1)
		for retry_count := 0; retry_count < retry_n; retry_count++ {
			err := extractPageInfoByAPI(key)
			if err == nil {
				if if_log {
					fmt.Printf("%s\n", strBGM(key, r.sep))
				}
				break
			} else {
				time.Sleep(150 * time.Millisecond)
			}
		}
	}
}

const (
	begin_number      int    = 28221450
	end_number        int    = 28222450
	default_delimiter string = ";"
	target_file_name  string = "output.csv"
	thread_number     int    = 20
	retry_number      int    = 5
	log_display       bool   = false
)

var (
	bgm_data_map = make(map[int]*string)
	bgm_mutex    = sync.Mutex{}
	begin_idx    int
	end_idx      int
	thread_n     int
	retry_n      int
	if_log       bool
)

func main() {
	begin_num := flag.Int("begin", begin_number, "起始 media ID")
	end_num := flag.Int("end", end_number, "终止 media ID")
	delimiter := flag.String("delimiter", default_delimiter, "分隔号")
	target_path := flag.String("output", target_file_name, "输出文件名")
	thread_num := flag.Int("thread", thread_number, "线程数")
	retry_num := flag.Int("retry", retry_number, "最大重试次数")
	logging := flag.Bool("log", log_display, "是否启用日志打印 (default false)")
	flag.Parse()

	begin_idx = *begin_num
	end_idx = *end_num
	thread_n = *thread_num
	retry_n = *retry_num
	if_log = *logging
	count := end_idx - begin_idx

	wg := &sync.WaitGroup{}
	queue := make([]*WorkRoutine, thread_n)

	for i := 0; i < thread_n; i++ {
		queue[i] = &WorkRoutine{idx: i}
		step := count / thread_n
		queue[i].begin = begin_idx + queue[i].idx*step
		queue[i].end = queue[i].begin + step
		if queue[i].end > end_idx {
			queue[i].end = end_idx
		}
	}

	defer saveData(*target_path, *delimiter)

	fmt.Println("开始扫描：")
	bar := progressbar.Default(
		int64(count),
		"BiliBGM",
	)

	start := time.Now()
	for i := 0; i < *thread_num; i++ {
		wg.Add(1)
		go queue[i].startRoutine(wg, bar)
	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("扫描用时：%s\n", elapsed)
}
