(function($, window, document) {
    //头部搜索栏的打开和隐藏
    $.fn.searchBar = function() {
        return this.each(function(i, v) {
            $(v).find('.search_btn').click(function() { //搜索
                var t = $(this),
                    p_b = $(this).parent();
                if ($(v).find('.sbi_inp').val() == '') {
                    if ($(v).css('width') == '35px') {
                        $(v).siblings('.navbar_ul').hide();
                        t.addClass('sb_bottomLine');
                        $(v).css('overflow', 'visible');
                        $(v).animate({
                            width: '420px'
                        }, 200);
                        $('.navbar_collapse').show();
                    } else {
                        $(v).animate({
                            width: '35px'
                        }, 200, 'linear', function() {
                            $(v).css('overflow', 'hidden');
                            $(v).siblings('.navbar_ul').show();
                            $('.navbar_collapse').hide();
                            t.removeClass('sb_bottomLine');
                        });
                    }
                }
            });

            var isClick = false;
            $(v).find('.c-search_sort_li').click(function() {
                isClick = true;
            });
            $(v).find('.sbi_inp').blur(function() {
                setTimeout(function() {
                    if (isClick) {
                        $(v).find('.sbi_inp').focus();
                        isClick = false;
                        return;
                    }
                    if ($(v).css('width') == '35px') {
                        return;
                    }
                    $(v).animate({
                        width: '35px'
                    }, 200, 'linear', function() {
                        $(v).css('overflow', 'hidden');
                        $(v).find('.search_btn').removeClass('sb_bottomLine');
                        $(v).siblings('.navbar_ul').show();
                        $('.navbar_collapse').hide();
                        $(v).find('.sbi_inp').val('');
                    });
                }, 100);
            });
            $(v).find('.sbs_ul').children('li').click(function() {
                $('#search_box .sbs_ul').hide();
                var txt = $(this).html(),
                    type = $(this).attr('data-type');
                $(v).find('#sbs_s').html(txt);
                $(v).find('#sbs_s').attr('data-type', type);
            });
        });
    }
    // 搜索功能
    $.fn.searchFun = function() {
        return this.each(function(i, v) {
            $(v).parents('.c-search_box').find('.c-search_sort_box').on('mouseover', function() {
                $(this).parents('.c-search_box').find('.c-search_sort_ul').show();
            });

            $(v).parents('.c-search_box').find('.c-search_sort_box').on('mouseleave', function() {
                $(this).parents('.c-search_box').find('.c-search_sort_ul').hide();
            });

            function moveEnd(obj) {
                obj.focus();
                var len = obj.value.length;
                if (document.selection) {
                    var sel = obj.createTextRange();
                    sel.moveStart('character', len);
                    sel.collapse();
                    sel.select();
                } else if (typeof obj.selectionStart == 'number' &&
                    typeof obj.selectionEnd == 'number') {
                    obj.selectionStart = obj.selectionEnd = len;
                }
            }
            $(v).parents('.c-search_box').find('.c-search_sort_li').on('click', function() {
                var txt = $(this).html(),
                    type = $(this).attr('data-type');
                $(this).parents('.c-search_box').find('.c-search_sort_default').html(txt).attr('data-type', type);
                $(this).parents('.c-search_box').find('.c-search_sort_ul').hide();
            });

            function searFun(obj, target) {
                var key = obj.parents('.c-search_box').find('.c-keywords').val(),
                    type = obj.parents('.c-search_box').find('.c-search_sort_default ').attr('data-type');
                key = key.replace(/(^\s*)|(\s*$)/g, '').replace(/^\s*$/, '');
                if (key == '') {
                    obj.parents('.c-search_box').find('.c-keywords').val('');
                } else {
                    obj.parents('.c-search_box').find('.c-keywords').val('');
                    if (target) {
                        window.open('http://www.cndesign.com/Query/To?type=' + type + '&key=' + escape(key));
                    } else {
                        window.location.href = 'http://www.cndesign.com/Query/To?type=' + type + '&key=' + escape(key);
                    }
                }
            }
            $(v).on('click', function() {
                if ($(this).parents('.c-search_box').hasClass('_self')) {
                    var target = 1;
                    searFun($(this), target);
                } else {
                    searFun($(this));
                }
            });
            // 回车检索
            $(v).parents('.c-search_box').find('.c-keywords').on('keydown', function(e) {
                if (e.keyCode == 13) {
                    if ($(this).parents('.c-search_box').hasClass('_self')) {
                        target = 1;
                        searFun($(this), target);
                    } else {
                        searFun($(this));
                    }
                }
            });
        });
    }
})(jQuery, window, document);

// 判断浏览器版本
var ie = -1;
var browser = navigator.appName;
var b_version = navigator.appVersion;
var version = b_version.split(";");
if (version.length > 1) {
    var trim_Version = version[1].replace(/[ ]/g, "");
    if (browser == "Microsoft Internet Explorer" && trim_Version == "MSIE8.0") {
        ie = 8;
    }
}

$(function() {
    //头部搜索栏 打开和隐藏
    $('#search_box').searchBar();

    //右上角消息和上传图标的动效
    $('.top_btns').children('a').mouseover(function(event) {
        event.stopPropagation();
        $(this).stop().animate({
            backgroundPositionY: '-60px'
        }, 200);
    });
    $('.top_btns').children('a').mouseout(function(event) {
        event.stopPropagation();
        $(this).stop().animate({
            backgroundPositionY: '0px'
        }, 200);
    });

    // 非首屏二级导航默认选中
    var index_sub_idx = $('.c-sub-nav_li > a.active');
    $('.c-sub-nav_li > a').on('mouseover', function() {
        $('.c-sub-nav_li > a').removeClass('active');
        $(this).addClass('active');
    });
    // index_sub_idx.attr('id','tar');
    $('.c-sub-nav_ul').on('mouseleave', function() {
        $('.c-sub-nav_li > a').removeClass('active');
        if (index_sub_idx.length > 0) {
            index_sub_idx.addClass('active');
        }
    });

    $('.head').css('width', $(document).width() + 'px');
    var docWidth = $(document.body).width();
    $(window).resize(function() {
        docWidth = $(document.body).width();
        $('.nb_div').css({
            'width': docWidth,
            'left': '0px'
        });
    });
    $(window).scroll(function() {
        if (docWidth <= 1200 && $('.head').css('position') == 'fixed') {
            $('.head').css('left', -$(window).scrollLeft() + 'px');
            $('.c-sub-nav').css('left', -$(window).scrollLeft() + 'px').css('min-width', '1200px');
            $('.nb_div').css('left', -$(window).scrollLeft() + 'px').css('width', $(document).width() + 'px');
        } else {
            $('.head').css('left', '0px');
            $('.c-sub-nav').css('left', '0');
        }
    });

    $('div[class*=-sub-nav]').hide();
});

jQuery.zqj = {
    switchNavFn: function() {
        $('#search_box').find('.search_btn').removeClass('sb_bottomLine');
        $('#search_box').css({
            'width': '35px',
            'overflow': 'hidden'
        });
        $('#search_box').find('.sbi_inp').val('');
        var top = $(window).scrollTop();
        if (top > 120) {
            $('.nav_logo>img').attr('src', PicHub + '/Content/images/logo_a.png');
            $('.nav_logo').attr('href', 'http://www.cndesign.com');
            if (ie !== 8) {
                $('.nav').css('margin-top', '60px');
            }
            $('.head').css('width', '100%').css('min-width', '1200px').css('position', 'fixed').css('top', '0px').css('z-index', '800');
            $('.head .navbar_ul').show();
            $('.head #search_box').css('display', 'block');
            $('.head .navbar_a.active').addClass('hover');
        } else {
            if (typeof logo == 'undefined') {
                $('.nav_logo>img').attr('src', PicHub + '/Content/images/logo_b.png');
                $('.nav_log').attr('href', window.location.origin);
            } else {
                $('.nav_logo>img').attr('src', logo);
                $('.nav_log').attr('href', window.location.origin);
            }
            if (ie !== 8) {
                $('.nav').css('margin-top', '0');
            }
            $('.head').css('position', 'static');
            $('.head .navbar_ul').hide();
            $('.head #search_box').css('display', 'none');
        }
    },
    //顶部导航条切换
    switchNav: function() {
        var me = this;
        return $(window).on('scroll', function() {
            me.switchNavFn();
        });
    }
}

$.zqj.switchNavFn(); // 
$.zqj.switchNav(); //顶部导航条切换

$('.c-search_btn').searchFun(); //调用搜索

/* 首屏漂浮二级导航 */
$(function() {
    $.zqj.defIdx = $('.nav .nav_ul_v1 .nav_a_v1').index($('.nav_a_v1.active')); //默认选中项
    nav();
    // 首屏二级导航鼠标移入移出状态切换
    var $def = $('.cnd-static .c-sub-nav_li a.active'); //默认选中元素
    $('.cnd-static').on('mouseover', '.c-sub-nav_li', function() {
        $(this).parents('.c-sub-nav_ul').find('a').removeClass('active');
        $(this).find('a').addClass('active');
    })
    $('.cnd-static').on('mouseout', '.c-sub-nav_li', function() {
        $(this).find('a').removeClass('active');
        $def.addClass('active');
    });
});

function nav() {
    if ($('.subNavWrapper').length > 0) { //首屏绝对定位二级导航
        $('.nav_a_v1').hover(function() {
            $.zqj.idx = $('.nav_a_v1').index(this);
            if ($(this).hasClass('haveSub')) {
                $('.subNavWrapper').show();
            }
            $('.nav_a_v1').removeClass('active');
            $(this).addClass('active');
            $('.subNavSet div').hide();
            $('.subNavSet div').eq($.zqj.idx).show();
        }, function() {
            $('.subNavWrapper').hide();
            $(this).removeClass('active');
            if ($.zqj.defIdx != -1) {
                $('.nav .nav_li_v1').find('.nav_a_v1').eq($.zqj.defIdx).addClass('active');
            }
        });
        $('.subNavSet div').hover(function() {
            $('.subNavWrapper').show();
            $(this).show();
            $('.nav .nav_li_v1 .nav_a_v1').removeClass('active');
            $('.nav .nav_li_v1').eq($.zqj.idx).find('.nav_a_v1').addClass('active');
            $('.nav .nav_li_v1').eq($.zqj.idx).siblings().find('.nav_a_v1').removeClass('active');
        }, function() {
            $('.subNavWrapper').hide();
            $('.nav .nav_li_v1 .nav_a_v1').removeClass('active');
            if ($.zqj.defIdx != -1) {
                $('.nav .nav_li_v1').find('.nav_a_v1').eq($.zqj.defIdx).addClass('active');
            }
        });
    }
    if ($('.cnd-static').length > 0) { //首屏正常文档流中二级导航
        $('.nav_a_v1').hover(function() {
            $.zqj.idx = $('.nav_a_v1').index(this);
            $('.nav_a_v1').removeClass('active');
            $(this).addClass('active');
            $('.zryd-place div').removeClass('show');
            $('.zryd-place').find('div').eq($.zqj.idx).addClass('show');
        });

        $(".zryd-nav").on('mouseleave', function() {
            $.zqj.idx = $.zqj.defIdx;
            $('.nav .nav_li_v1 .nav_a_v1').removeClass('active');
            $('.nav_a_v1').eq($.zqj.idx).addClass('active');
            $('.zryd-place').find('div').removeClass('show');
            $('.zryd-place').find('div').eq($.zqj.idx).addClass('show');
        });
    }
}
/* 首屏漂浮二级导航 */