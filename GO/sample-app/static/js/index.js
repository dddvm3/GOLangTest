const root = document.documentElement;

var reqestURL = "/static/Json/QuestionList.json"

function rectime(){
    var finishtime= Date.now;
    var time = Date.now-finishtime;
    console.log(time);
};
    
$(function(){
    $(".questioncontent").ready(function(){showcontenttext(this)});
    $(".CountBar").ready(function(){countBar(this)});
})

function showcontenttext(element){
    var text = $("#content").attr("text");
    var tl = gsap.timeline();tl.to("tspan",{duration:2,text:text,ease:"none"});
}

function countBar(element){
    var count=$("#countbar").attr("count");
    console.log(count);
    root.style.setProperty('--innerwidth',count/10);
}

var postdata;
var answer;
var userid;
var postanswer;
var quesstionCount;

function checkanswer(element){
    answer=$("#checkanswer").attr("answer");
    userid=$("#userID").attr("userid");
    postanswer = document.getElementById('postanswer');
    postanswer.setAttribute("value",element);
    quesstionCount = $("#questionCount").attr("questionCount");
    
    postdata = {userID:userid,answer:element}

    console.log(postdata);
    console.log(answer);
    console.log(userid);
    console.log(element);
    console.log(quesstionCount);
    $("#checkanswer").show();

    $("#answerA").html(answer=="A"?"O":"X");
    $("#answerB").html(answer=="B"?"O":"X");
    $("#answerC").html(answer=="C"?"O":"X");
    $("#answerD").html(answer=="D"?"O":"X");

    $("#countdown").attr("src","/static/Images/countdown_3.png");

    count_index=3;
    countdown();
 
}

function countdown(){
    setTimeout(() => {
        count_index--;
        if(count_index>0){
            $("#countdown").attr("src","/static/Images/countdown_"+count_index+".png");
            countdown();
        }
        else{
            if(quesstionCount!=10)
            {
                sendPost();
            }
        }
    }, 1000);
}

function sendPost(){
    const jsonData = JSON.stringify(postdata);

    fetch('http://10.10.2.40:8000/QuestionPage',{
        method:'POST',
        headers: {
            'Content-Type': 'application/json'
          },
        body:jsonData
    })
    .then(response=>{
        return response.text();
    }).then(html => {
        
        const newDoc = new DOMParser().parseFromString(html,'text/html');

        window.location.reload();

    }).catch(error => {
        console.error('There was a problem with the fetch operation:', error);
    });
}