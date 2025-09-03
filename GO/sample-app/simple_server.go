package main

import (
	// "encoding/json"

	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
)

var questionList QuestionList
var usequestionList []UsequestionList

type Message struct {
	UserID string `json:"userID"`
	Answer string `json:"answer"`
}

func main() {
	ReadJson()
	fmt.Println(len(questionList.Questions))

	router := gin.Default()
	router.Static("/static", "./static")

	router.LoadHTMLGlob("templates/*")

	router.GET("/home", func(c *gin.Context) {
		node, err := snowflake.NewNode(64)
		if err != nil {
			fmt.Println(err)
			return
		}
		userID := node.Generate().Base36()
		// fmt.Println(userID)
		usequestionList = append(usequestionList, UsequestionList{userID: userID, QuestionNum: generateRandomNumbder(len(questionList.Questions), 10)})

		c.HTML(200, "index.html",
			gin.H{"title": "歡迎光臨 大問答",
				"userID":  userID,
				"Contant": "你是最後一位參賽者"})

	})

	router.GET("/static", func(ctx *gin.Context) {
		ctx.XML(200, "index.css")
	})

	router.POST("/MainPage", func(ctx *gin.Context) {

		UserID := ctx.PostForm("userID")
		// fmt.Println(UserID)
		var Num int
		for i := 0; i < len(usequestionList); i++ {
			if usequestionList[i].userID == UserID {
				Num = i
				break
			}
		}

		questionNum := len(usequestionList[Num].Succse)

		fmt.Println(usequestionList)
		if len(usequestionList[Num].Succse) != 10 {
			ctx.HTML(200, "questionPage.html",
				gin.H{"questionNum": strconv.FormatInt(int64(questionNum+1), 10),
					"title":    "第" + strconv.FormatInt(int64(questionNum+1), 10) + "題",
					"userID":   UserID,
					"Question": questionList.Questions[usequestionList[Num].QuestionNum[questionNum]].Content,
					"OptionA":  questionList.Questions[usequestionList[Num].QuestionNum[questionNum]].A,
					"OptionB":  questionList.Questions[usequestionList[Num].QuestionNum[questionNum]].B,
					"OptionC":  questionList.Questions[usequestionList[Num].QuestionNum[questionNum]].C,
					"OptionD":  questionList.Questions[usequestionList[Num].QuestionNum[questionNum]].D,
					"Answer":   questionList.Questions[usequestionList[Num].QuestionNum[questionNum]].Answer,
				})
		} else {

			ctx.HTML(200, "EndPage.html",
				gin.H{"title": "結果",
					"NumberofAnswers": strconv.FormatInt(int64(usequestionList[Num].Score), 10)})
			fmt.Println(usequestionList)
		}
	})

	router.POST("/QuestionPage", func(ctx *gin.Context) {
		// fmt.Println(ctx.Request.Body)
		var message Message
		if er := ctx.BindJSON(&message); er != nil {
			ctx.JSON(400, gin.H{"error": er.Error()})
			return
		}

		UserID := message.UserID
		fmt.Println(UserID)

		var Num int
		for i := 0; i < len(usequestionList); i++ {
			if usequestionList[i].userID == UserID {
				Num = i
				break
			}
		}
		answer := message.Answer

		checkquestionNum := len(usequestionList[Num].Succse)
		if questionList.Questions[usequestionList[Num].QuestionNum[checkquestionNum]].Answer == answer {
			usequestionList[Num].Succse = append(usequestionList[Num].Succse, true)
			usequestionList[Num].Score = usequestionList[Num].Score + 1
		} else {
			usequestionList[Num].Succse = append(usequestionList[Num].Succse, false)
		}

		questionNum := len(usequestionList[Num].Succse)

		// fmt.Println(answer)
		fmt.Println(usequestionList)

		if len(usequestionList[Num].Succse) != 10 {
			ctx.HTML(200, "questionPage.html",
				gin.H{"questionNum": strconv.FormatInt(int64(questionNum+1), 10),
					"title":    "第" + strconv.FormatInt(int64(questionNum+1), 10) + "題",
					"userID":   UserID,
					"Question": questionList.Questions[usequestionList[Num].QuestionNum[questionNum]].Content,
					"AnswerA":  questionList.Questions[usequestionList[Num].QuestionNum[questionNum]].A,
					"AnswerB":  questionList.Questions[usequestionList[Num].QuestionNum[questionNum]].B,
					"AnswerC":  questionList.Questions[usequestionList[Num].QuestionNum[questionNum]].C,
					"AnswerD":  questionList.Questions[usequestionList[Num].QuestionNum[questionNum]].D})
		} else {

			ctx.HTML(200, "EndPage.html",
				gin.H{"title": "結果",
					"NumberofAnswers": strconv.FormatInt(int64(usequestionList[Num].Score), 10)})
			fmt.Println(usequestionList)
		}
	})

	// router.POST("/EndPage", func(ctx *gin.Context) {
	// 	ctx.HTML(200, "EndPage.html", gin.H{})
	// })

	router.Run(":8000")

}

func generateRandomNumbder(end int, count int) []int {
	if end < 0 || end < count {
		return nil
	}

	nums := make([]int, 0)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for len(nums) < count {
		num := r.Intn(end)

		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}

		if !exist {
			nums = append(nums, num)
		}
	}
	return nums
}

func ReadJson() {

	questionFile := "./static/Json/QuestionList.json"
	fl, err := os.Open(questionFile)
	if err != nil {
		fmt.Println(questionFile)
	}
	defer fl.Close()
	buf := make([]byte, 8192)
	for {
		n, _ := fl.Read(buf)
		if n == 0 {
			break
		}
		// os.Stdout.Write(buf[:n])
		fmt.Println(n)
		str := string(buf[:n])

		fmt.Println(str)
		json.Unmarshal([]byte(str), &questionList)
	}
}

type Aquestion struct {
	Content string
	A       string
	B       string
	C       string
	D       string
	Answer  string
}

type QuestionList struct {
	Questions []Aquestion
}

type Usequestion struct {
	QuestionNum int
	Succse      bool
}

type UsequestionList struct {
	userID      string
	QuestionNum []int
	Succse      []bool
	Score       int
}
