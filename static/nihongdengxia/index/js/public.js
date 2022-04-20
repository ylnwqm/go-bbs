layui.use(['element', 'form', 'layer', 'carousel'], function(){ 
var element = layui.element
,form = layui.form
,layer = layui.layer
,carousel = layui.carousel;

carousel.render({
    elem: '#index'
    ,width: '340px' 
    ,height: '135px' 
    ,arrow: 'hover'
    ,anim: 'default'
  });
});



$(function() {
	$('.article-body').find('img').each(function() {
		var _this = $(this);
		var _src = _this.attr("src");
		var _alt = _this.attr("alt");
		_this.wrap('<a data-fancybox="images" href="' + _src + '" data-caption="' + _alt + '"></a>');
	})
})



function opensearcbox() {
	var search = $("#search");
	layui.use(['layer', 'form'], function() {
		var search_index = layer.open({
			type: 1,
			skin: "searchskin",
			content: search,
			btn: 0,
			close: 0,
			title: false,
			area: ["678px", "390px"],
			end: function() {
				$("#searchBox").hide();
			}
		});
		$("#searchBox").show();
		$(".closesearch").click(function(e) {
			layer.close(search_index);
		});
	});
}


// if ($(".backtop").length > 0) {
// 	backtopS();
// 	$(window).scroll(function(e) {
// 		backtopS();
// 	});

// 	function backtopS() {
// 		var scroH = $(this).scrollTop();
// 		var footerHeight = 0;
// 		if ($('.footer').length > 0) {
// 			mTop = $('.footer')[0].offsetTop;
// 			footerHeight = footerHeight + $(".footer").outerHeight();
// 		} else {
// 			mTop = $('.footer')[0].offsetTop;
// 		}
// 		footerHeight = $(".footer").outerHeight() - parseInt($("footer").css("margin-top")) - parseInt($("footer").css("margin-top")) - parseInt($("footer").css("margin-top"));
// 		sTop = $(window).scrollTop();
// 		result = mTop - sTop - parseInt($("footer").css("margin-top"));
// 		if (scroH > 200) {
// 			$(".backtop").fadeIn(400);
// 			if (scroH >= $("body").height() - $(window).height() - footerHeight) {
// 				$(".backtop").css("bottom", $(window).height() - result);
// 			} else {
// 				$(".backtop").css("bottom", "");
// 			}
// 		} else {
// 			$(".backtop").fadeOut(400);
// 		}
// 	}
// }
// $(".backtop").click(function(e) {
// 	$('body,html').animate({
// 		scrollTop: 0
// 	}, 800);
// });



function ffix() {
	var u;
	if (floatFix) {
		var n = $(".right.public .box"),
			t = $(window).scrollTop(),
			f = n.offset().top - t,
			i = n.height(),
			r = $(window).height();
		i + 40 > r ? floatTop < t - (i + 40 - r) ? n.hasClass("fix-bottom") || n.addClass("fix-bottom") : n.hasClass("fix-bottom") && n.removeClass("fix-bottom") : floatTop < 60 + t ? n.hasClass("fix-top") || n.addClass("fix-top") : n.hasClass("fix-top") && n.removeClass("fix-top");
		u = $(window).height() - $(".footer").offset().top + $(window).scrollTop() + 20;
		u > 0 ? n.hasClass("at-bottom") || n.addClass("at-bottom") : n.hasClass("at-bottom") && n.removeClass("at-bottom");
		lastScrollTop = t
	}
}


    // r-list 切换
    $(".r-list .bar li").hover(function () {
        var tabList = $(this).parent().parent();
        tabList.find(".bar li").removeClass("sel");
        //$("#rank .bar li").removeClass("sel");
        $(this).addClass("sel");

        tabList.find(".bd").removeClass("sel");
        //$("#rank .bd").removeClass("sel");
        //$("#rank .cm").removeClass("sel");
        var id = "#d-" + $(this).data("id");
        tabList.find(id).addClass("sel");
    });
$(".article-tit .info span:last-child").click(function() {
	var target_top = $(".article-pinglun").offset().top; 
    $("html,body").animate({scrollTop: target_top}, 1000); 
});   

$('.site-menu a').click(function() {
    $(this).addClass('hover').siblings().removeClass('hover');
	var href = $(this).attr("href");
	var pos = $(href).offset().top - 135;
	$("html,body").animate({
		scrollTop: pos
	}, 400);
	return false;
});


$(function() {
	$("header .user-open").click(function() {
		$("header .user-open .user-open-info").toggleClass("hover");
		$("header .user-open .tx").toggleClass("hover");
	})
})
