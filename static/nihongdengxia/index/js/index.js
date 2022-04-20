

$(document).ready(function() {
	var NLNUM = $('.index-news-list .page a');
	NLNUM.click(function() {
		NLNUM.removeClass('current');
		$(this).addClass('current');
		var NUMB = $(this).index() + 1;
		if ($(this).attr('class') == 'current') {
			$('.block').css({
				'display': 'none'
			});
			$('.new-list-' + String(NUMB)).css({
				'display': 'block'
			})
			$(window).trigger('scroll');
		} else {
			$('.new-list .block').css({
				'display': 'none'
			})
		}
	});
});

$('.list-article-tab > .type > a').on('mouseover', function() {
	$(this).addClass('hover').siblings().removeClass('hover');
	$(this).parents('.list-article-tab').find('.plate-list').eq($(this).index()).addClass('ing').siblings('.plate-list').removeClass('ing');
	$(window).trigger('scroll');
});
$('.isoft .list-soft-tab > .type > a,.igame .list-soft-tab > .type > a').on('mouseover', function() {
	$(this).addClass('hover').siblings().removeClass('hover');
	$(this).parents('.list-soft-tab').find('.plate-list').eq($(this).index()).addClass('ing').siblings('.plate-list').removeClass('ing');
	$(window).trigger('scroll');
});
$('.list-course-tab .type > a').on('mouseover', function() {
	$(this).addClass('hover').siblings().removeClass('hover');
	$(this).parents('.list-course-tab').find('.plate-list').eq($(this).index()).addClass('ing').siblings('.plate-list').removeClass('ing');
	$(window).trigger('scroll');
});

$(".hot-soft .type a").hover(function(e) {
	$(".hot-soft .type .hover").removeClass("hover");
	$(this).addClass("hover"); /*计算当前应该偏移的高度*/
	var index = $(".hot-soft .type .hover").index();
	var x = $(".hot-soft").width() * index + parseInt($(".hot-soft").css("margin-right")) * index;
	$(".hot-soft .soft-plate").css("transform", "translateX(-" + x + "px)");
});
$(".hot-soft .plate-list ul").each(function() {
	$(this).find("li:first").addClass("now"), $(this).find("li").hover(function() {
		$(this).addClass("now").siblings().removeClass("now")
	})
})
$(".list-soft-tab .type,.list-article-tab .type,.list-course-tab .type").each(function() {
		$(this).find("a:first").addClass("hover");
})
$(".isoft .list-soft-plate, .igame .list-soft-plate,.list-article-tab .list-article-plate,.list-course-tab .list-course-plate").each(function() {
		$(this).find(".plate-list:first").addClass("ing");
})

layui.use(['laytpl', 'util', 'layer'], function () {

	$(document).click(function (e) {
        $(".selectul").hide();
    }); /*首页推荐栏悬浮后展开下拉框*/
    $(".indexWebList li a .xl").hover(function () {
        $(this).parents("li").find(".info").show();
        $(this).parents("li").find(".xl").addClass("xlhover")
    });
    $(".indexWebList li").mouseleave(function () {
        $(this).find(".info").hide();
        $(this).find(".xl").removeClass("xlhover")
    });
    $(".indexWebList li a .xl").click(function (e) {
        e.preventDefault();
    });
});

$(".index-news .index-new-list").each(function() {
		$(this).find(".index-news-list:first").removeClass('layui-hide');
})
$('.index-news .news-type a').on('click', function() {
	$(this).addClass('hover').siblings().removeClass('hover');
	$('.index-new-list').find('.index-news-list').eq($(this).index()).removeClass('layui-hide').siblings('.index-news-list').addClass('layui-hide');
	$(window).trigger('scroll');
});


$(".index-news .index-news-list").each(function() {
		$(this).find(".new-list:first").removeClass('layui-hide');
})

$('.index-news-list .page a').on('click', function() {
	
	$(this).addClass('current').siblings().removeClass('current');
	$(this).parents('.index-news-list').find('.new-list').eq($(this).index()).removeClass('layui-hide').siblings('.new-list').addClass('layui-hide');
	$(window).trigger('scroll');
});