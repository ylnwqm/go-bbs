/* 轮播 */
var curIndex = 0, //当前index
    imgLen = $(".imgList li").length; //图片总数
// 定时器自动变换5秒每次
var autoChange = setInterval(function(){
    if(curIndex < imgLen-1){
        curIndex ++;
    }else{
        curIndex = 0;
    }
    //调用变换处理函数
    changeTo(curIndex);
},5000);
/*
 //左箭头滑入滑出事件处理
 $("#prev").hover(function(){
    //滑入清除定时器
    clearInterval(autoChange);
 },function(){
    //滑出则重置定时器
    autoChangeAgain();
 });
//左箭头点击处理
 $("#prev").click(function(){
    //根据curIndex进行上一个图片处理
    curIndex = (curIndex > 0) ? (--curIndex) : (imgLen - 1);
    changeTo(curIndex);
 });
//右箭头滑入滑出事件处理
 $("#next").hover(function(){
    //滑入清除定时器
    clearInterval(autoChange);
 },function(){
    //滑出则重置定时器
    autoChangeAgain();
 });
 //右箭头点击处理
 $("#next").click(function(){
    curIndex = (curIndex < imgLen - 1) ? (++curIndex) : 0;
    changeTo(curIndex);
 });
 */
 $('.imgList').on('mouseover','li>a',function () {
    clearInterval(autoChange);
 })
 $('.imgList').on('mouseout','li>a',function () {
    autoChangeAgain();
 })
 /*
 $('.bg').on({
    'mouseover':function () {
        clearInterval(autoChange);
    },
    'mouseout':function () {
        autoChangeAgain();
    }
 })
 $('.infoList').on('mouseover','.infoOn',function () {
    clearInterval(autoChange);
 })
 $('.infoList').on('mouseout','.infoOn',function () {
    autoChangeAgain();
 })
 */
//对右下角按钮index进行事件绑定处理等
$(".indexList").find("li").each(function(item){
    $(this).hover(function(){
        clearInterval(autoChange);
        changeTo(item);
        curIndex = item;
    },function(){
        autoChangeAgain();
    });
});
//清除定时器时候的重置定时器--封装
function autoChangeAgain(){
    autoChange = setInterval(function(){
        if(curIndex < imgLen-1){
            curIndex ++;
        }else{
            curIndex = 0;
        }
        //调用变换处理函数
        changeTo(curIndex);
    },5000);
}
function changeTo(num){
    var goLeft = num * 760;
    $(".imgList").stop().animate({left: "-" + goLeft + "px"},400);
    $(".infoList").find("li").removeClass("infoOn").eq(num).addClass("infoOn");
    $(".indexList").find("li").removeClass("indexOn").eq(num).addClass("indexOn");
}
/* 轮播 End */

$(function(){

    /* 文章切换 */
    $('.articles .links_ul .links_a').on('mouseover',function(){
        var blockNum = $('.articles .links_ul .links_a').index($(this)),
            href = $(this).attr('href');
        $(this).parents('.articles').find('.more').attr('href',href); //替换更多的超链接
        if($('.articles .articles_ul').eq(blockNum).html()!=''){
            $('.articles .articles_ul').hide();
            $('.articles .articles_ul').eq(blockNum).show();
            $('.articles .links_ul .links_a').removeClass('active');
            $(this).addClass('active');
            return;
        }else if($(this).hasClass('active')){
            return;
        }else {
            $('.articles .links_ul .links_a').removeClass('active');
            $(this).addClass('active');
            var lang = $(this).parent().attr('lang');
            $.get(lang, function (res) {
                obj = res;
                var $ul = $('<ul class="articles_ul clear"></ul>');
                for (var i = 0, str = ''; i < obj.length; i++) {
                    str += '<li class="articles_li lf clear">' +
                            '<a class="articles_cover_box lf" href="' + obj[i].newsUrl + '" target="_blank">' +
                            '<img class="articles_cover zoom" src="' + obj[i].cover + '" alt="">' +
                            '</a>' +
                            '<div class="articles_txt">' +
                            '<div class="sm_title_box">' +
                            '<a href="' + obj[i].classInfo[0].classUrl + '" rel="nofollow">'+ obj[i].classInfo[0].className +'</a>' +
                            ' | ' +
                            '<a href="' + obj[i].newsUrl + '" target="_blank">' + obj[i].title + '</a>' +
                            '</div>' +
                                '<p class="articles_brief">' +
                                    '<a href="'+ obj[i].newsUrl +'" target="_blank">'+ obj[i].introduce +'</a>' +
                                '</p>' +
                            '</div>' +
                        '</li>';
                }
                $ul.html(str);
                $('.articles .content_box .articles_ul').eq(blockNum).replaceWith($ul);
                $('.articles .articles_ul').hide();
                $('.articles .articles_ul').eq(blockNum).show();
                $('.zoom').zoomImgRollover();
            });
        }
    })
    /* 文章切换 End */

    $('.works').on('mouseover','.pi_person',function () {   //当鼠标移入作品列表中作者头像时
        $(this).parents('.prod_info_box').find('.pi_name').addClass('hover');   //修改对应作者姓名字体颜色
    })
    $('.works').on('mouseout','.pi_person',function () {    //当鼠标移出作品列表中作者头像时
        $(this).parents('.prod_info_box').find('.pi_name').removeClass('hover');    //还原对应作者姓名字体颜色
    })
    $('.design_li').on('mouseover','.avatar',function(){    //当鼠标移入设计师、设计公司列表中头像时
        $(this).parents('.design_li').find('.design_name').addClass('hover');   //修改对应设计师、设计公司名称字体颜色
    })
    $('.design_li').on('mouseout','.avatar',function(){ //当鼠标移出设计师、设计公司列表中头像时
        $(this).parents('.design_li').find('.design_name').removeClass('hover');    //还原对应设计师、设计公司名称字体颜色
    })
})