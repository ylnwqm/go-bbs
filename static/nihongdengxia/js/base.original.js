(function ($, window, document) {
    $.picUrl = '/HomeAsync/CropImage', $.upUrl = '/HomeAsync/upload'
    var curPage = 0, pages = 0;
    $.jcropApi = null;
    $.extend({
        check: function (val, type) {
            var r;
            switch (type) {
                case 'account':
                    r = /(^[\w.\-]+@(?:[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*\.)+[A-Za-z]{2,6}$)|(^1\d{10}$)|(^\d{3,}$)/.test(val);
                    break;
                case 'pwd':

            }
            return r;
        },
        showmsg: function (msg) {
            var timer = null;
            var w_m = $('<div class="w_p"><div class="w_bg"></div><div class="w_t">' + msg + '</div></div>');
            w_m.appendTo('body').css({ left: Math.floor(($(window).width() - $('.w_p').innerWidth()) / 2) + 'px', top: $(window).scrollTop() + Math.floor(($(window).height() - $('.w_p').innerHeight()) / 2) + 'px' }).animate({ opacity: 1 }, 1000);
            timer = setTimeout(function () { w_m.animate({ opacity: 0 }, 1000, function () { $(this).detach() }) }, 1500);
        },
        checkModify: function (n, old) {
            var o = eval(n);
            var _o = eval(old);
            if (o.length !== _o.length) {
                return true;
            }
            else {
                for (i = 0; i < o.length; i++) {
                    if (o[i].name == _o[i].name && o[i].value != _o[i].value || o[i].name != _o[i].name) {
                        return true;
                    }
                }
            }
            this.showmsg("您还没有做任何修改！");
            return false;
        },
        lightbox: function (ops) {
            return '<div class="cm">'+
                        '<div class="lightBox"></div>'+
                        '<div class="f_p">'+
                            '<div class="op_title">'+
                                '<span>' + ops.title + '</span>'+
                                '<div class="close c"></div>'+
                            '</div>'+
                            '<div class="tx_f_pan"></div>'+
                            '<div class="pforms">'+
                                '<div class="form_group">'+
                                    '<div class="mb20">'+
                                        '<label class="form_m">' + ops.content + '</label>'+
                                    '</div>'+
                                '<div class="op_f">'+
                                    '<input type="button" class="sub btn btn_2" value="' + ops.btn + '">'+
                                    '<input type="button" class="btn default_btn btn_2 c" value="取消">'+
                                '</div>'+
                            '</div>'+
                        '</div>'+
                    '</div>';
        },
        preview: function(coords){
            // 发送给后台的数据
            $('#photo').data('x', coords.x);
            $('#photo').data('y', coords.y);
            $('#photo').data('w', coords.w);
            $('#photo').data('h', coords.h);
            var zoom_scale = 130 / coords.w,
                zoom_scale1 = 70 / coords.w,
                zoom_scale2 = 40 / coords.w;
            $('.tx_img130 img').css({
                width: Math.round(zoom_scale * $.jcropApi.getBounds()[0]) + 'px',
                height: Math.round(zoom_scale * $.jcropApi.getBounds()[1]) + 'px',
                marginLeft: '-' + Math.round(zoom_scale * coords.x) + 'px',
                marginTop: '-' + Math.round(zoom_scale * coords.y) + 'px'
            });
            $('.tx_img70 img').css({
                width: Math.round(zoom_scale1 * $.jcropApi.getBounds()[0]) + 'px',
                height: Math.round(zoom_scale1 * $.jcropApi.getBounds()[1]) + 'px',
                marginLeft: '-' + Math.round(zoom_scale1 * coords.x) + 'px',
                marginTop: '-' + Math.round(zoom_scale1 * coords.y) + 'px'
            });
            $('.tx_img40 img').css({
                width: Math.round(zoom_scale2 * $.jcropApi.getBounds()[0]) + 'px',
                height: Math.round(zoom_scale2 * $.jcropApi.getBounds()[1]) + 'px',
                marginLeft: '-' + Math.round(zoom_scale2 * coords.x) + 'px',
                marginTop: '-' + Math.round(zoom_scale2 * coords.y) + 'px'
            });
        },
        file_up: function (ops) {
            try {
                var uploader = new qq.FileUploader({
                    element: document.getElementById(ops.elem),
                    action: $.upUrl,
                    dragText: '',
                    multiple: ops.multiple,
                    params: { t: ops.type || 0 },
                    allowedExtensions: ['jpg', 'jpeg', 'png', 'gif'],
                    customHeaders: { 'Content-Type': 'multipart/form-data' },
                    acceptFiles: 'image/jpg,image/jpeg,image/png,image/gif',
                    sizeLimit: 10485760,
                    listElement: '',
                    uploadButtonText: ops.btnText,
                    debug: true,
                    inputName: 'file',
                    onSubmit: function (id, filename) {
                        ops.sub && ops.sub(id, filename);
                    },
                    onProgress: function (id, filename, loaded, total) {
                        if (typeof ops.pro === 'function') {
                            ops.pro(id, filename, loaded, total); return;
                        }
                        ops.root.find('#photo').css('background', '#333 url(../Content/images/loading.gif) no-repeat center center');
                    },
                    onComplete: function (id, fileName, responseJSON) {
                        /** responseJSON path: 图片src路径 width: 360 height: 图片高度
                        *   所有上传的图片都会处理成宽度360px（等比缩放），然后返回。
                        **/
                        if (responseJSON.Status != 0) {
                            $('#loading_img').hide();
                            if(typeof ops.com === 'function') {
                                ops.com(id, fileName, responseJSON);
                                return;
                            }
                            $.showmsg(responseJSON.Msg);
                            return false;
                        }else{
                            if (typeof ops.com === 'function') {
                                ops.com(id, fileName, responseJSON);
                                return;
                            }
                            var o_img = ops.root.find('#photo'); // 页面上裁剪图片
                            o_img.data(responseJSON.Data); // 将图片信息存储到$('#photo')
                            ops.root.removeClass('uped').addClass('edit_tx');

                            if($.jcropApi != null){$.jcropApi.destroy()}; // 销毁裁剪实例，为二次初始化裁剪插件准备
                            ops.root.find('#photo').attr({ 'src': o_img.data('path') }); // 修改待裁剪图片src
                            $('#loading_img').hide();
                            ops.root.find('.tx_pre img').attr('src', o_img.data('path')); // 修改右侧预览图片src

                            ops.root.find('#photo').Jcrop({
                                allowSelect: true,
                                aspectRatio: 1,
                                onChange: $.preview,
                                onSelect: $.preview,
                                boxWidth: 360,
                                boxHeight: 360,
                                minSize: 10
                            }, function(){
                                $.jcropApi = null;
                                $.jcropApi = this;
                                var w, h;
                                this.setImage(o_img.data('path'), function(){
                                    w = $.jcropApi.getBounds()[0];
                                    h = $.jcropApi.getBounds()[1];
                                    var x = w / 2,
                                        y = h / 2,
                                        x_ = 130 / 2 * $.jcropApi.getScaleFactor()[0],
                                        y_ = 130 / 2 * $.jcropApi.getScaleFactor()[1];
                                    $.jcropApi.setSelect([x - x_, y - y_, x + x_, y + y_]);
                                });
                            });
                            o_img.parent().bind('dblclick', function () {
                                $.ajax({
                                    url: $.picUrl,
                                    type: 'POST',
                                    data: { x: o_img.data('x'), y: o_img.data('y'), w: o_img.data('w'), h: o_img.data('h'), path: o_img.data('path'), t: 0 },
                                    success: function (data) {
                                        $('.cm').remove();
                                        var Data = $.parseJSON(data);
                                        if (Data.Status !== 0)
                                            return false;
                                        $('.cur_tx img').attr('src', Data.Data.LargeHead);
                                        $.showmsg('修改头像成功！');
                                    },
                                    error: function () {
                                        $.showmsg("服务器错误！");
                                    }
                                });
                            })
                        }
                    },
                    onError: function (id, fileName, xhr) {
                        if (typeof ops.error === 'function') {
                            ops.error(id, fileName, xhr);
                        }
                        $.showmsg('服务器繁忙，请稍后再试');
                    },
                    showMessage: function (message) { $.showmsg(message) }
                })
            } catch (e) {

            }
        },
        // 身份认证页面的上传图片插件
        file_up_authentication: function(ops) {
            function fnCom(id, fileName, responseJSON) {
                if (responseJSON) {
                    var d = responseJSON.Data;
                    var cur = $('#'+ops.elem).siblings('.dd').find('.i_p[data-id=' + id + ']');
                    if(responseJSON.Status == 0) {
                        var o_img = cur.find('.img').attr('src', d.path).data(d);
                        cur.find('.del').bind('click', function () {
                            cur.remove();
                        });
                    }else{
                        cur.remove();
                        $.showmsg(responseJSON.Msg);
                        return false;
                    }
                }
            }
            try {
                var uploader = new qq.FileUploader({
                    element: document.getElementById(ops.elem),
                    action: $.upUrl,
                    dragText: '',
                    multiple: false,
                    params: { t: 2 },
                    allowedExtensions: ['jpg', 'jpeg', 'png', 'gif'],
                    customHeaders: { 'Content-Type': 'multipart/form-data' },
                    acceptFiles: 'image/jpg,image/jpeg,image/png,image/gif',
                    sizeLimit: 10485760,
                    listElement: '',
                    uploadButtonText: ops.btnText || '',
                    debug: true,
                    inputName: 'file',
                    onSubmit: function (id, filename) {
                        var t = $('#'+ops.elem).siblings('.dd').children('.dd-list').html('<li class="i_p dd-item" data-id="' + id + '"><div class="i_s_v" title="点击图片设置封面"><img class="img crop" src="'+PicHub+'/Content/images/loaded.gif"></div><em class="del qq-upload-cancel" title="删除">删除</em></div></li>');
                        t.find('.del').bind('click', function () {
                            t.remove();
                            return false;
                        });
                    },
                    onProgress: function (id, filename, loaded, total) {
                        var t = $('#'+ops.elem).siblings('.dd').find('.i_p[data-id=' + id + ']');
                        t.find('.bar').css('width', (loaded / total) * 100 + '%').find('.del').bind('click', function () { t.remove() });
                        return;
                    },
                    onComplete: function (id, fileName, responseJSON) {
                        if (responseJSON.Status != 0) {
                            $('#loading_img').hide();
                            fnCom(id, fileName, responseJSON);
                            return;
                        }else{
                            fnCom(id, fileName, responseJSON);
                            return;
                        }
                    },
                    onError: function (id, fileName, xhr) {
                        var t = $('#'+ops.elem).siblings('.dd').find('.i_p[data-id=' + id + ']');
                        t.remove();
                        $.showmsg('服务器繁忙，请稍后再试');
                    },
                    showMessage: function (message) { $.showmsg(message) }
                })
            } catch (e) {

            }
        },
        // 联系方式-二维码  上传二维码
        upload_qrCode: function(ops) {
            function fnCom(id, fileName, responseJSON) {
                if (responseJSON) {
                    var d = responseJSON.Data;
                    var cur = $('#'+ops.elem).siblings('.dd').find('.i_p[data-id=' + id + ']');
                    if(responseJSON.Status == 0) {
                        var o_img = cur.find('.img').attr('src', d.path).data(d);
                        $('#qrCode').val(d.path);
                        cur.find('.del').bind('click', function () {
                            $('#qrCode').val('');
                            cur.remove();
                        });
                    }else{
                        cur.remove();
                        $.showmsg(responseJSON.Msg);
                        return false;
                    }
                }
            }
            try {
                var uploader = new qq.FileUploader({
                    element: document.getElementById(ops.elem),
                    action: $.upUrl,
                    dragText: '',
                    multiple: false,
                    params: { t: 2 },
                    allowedExtensions: ['jpg', 'jpeg', 'png', 'gif'],
                    customHeaders: { 'Content-Type': 'multipart/form-data' },
                    acceptFiles: 'image/jpg,image/jpeg,image/png,image/gif',
                    sizeLimit: 10485760,
                    listElement: '',
                    uploadButtonText: ops.btnText,
                    debug: true,
                    inputName: 'file',
                    onSubmit: function (id, filename) {
                        var t = $('#'+ops.elem).siblings('.dd').children('.dd-list').html('<li class="i_p dd-item" data-id="' + id + '"><div class="i_s_v"><img class="img crop" src="'+PicHub+'/Content/images/loaded.gif"></div><em class="del qq-upload-cancel" title="删除">删除</em></div></li>');
                        t.find('.del').bind('click', function () {
                            $('#qrCode').val('');
                            t.remove();
                            return false;
                        });
                    },
                    onProgress: function (id, filename, loaded, total) {
                        var t = $('#'+ops.elem).siblings('.dd').find('.i_p[data-id=' + id + ']');
                        t.find('.bar').css('width', (loaded / total) * 100 + '%').find('.del').bind('click', function () { t.remove() });
                        return;
                    },
                    onComplete: function (id, fileName, responseJSON) {
                        if (responseJSON.Status != 0) {
                            $('#loading_img').hide();
                            fnCom(id, fileName, responseJSON);
                            return;
                        }else{
                            fnCom(id, fileName, responseJSON);
                            return;
                        }
                    },
                    onError: function (id, fileName, xhr) {
                        var t = $('#'+ops.elem).siblings('.dd').find('.i_p[data-id=' + id + ']');
                        t.remove();
                        $.showmsg('服务器繁忙，请稍后再试');
                    },
                    showMessage: function (message) { $.showmsg(message) }
                })
            } catch (e) {

            }
        },
        blackHandler: function (page) {
            var pages = 1;
            $.get('/HomeAsync/GetBlackList', { page: page }, function (data) {
                var temp = template('blist_temp', data);
                if (data.list.length == 0) {
                    $('.state_pass').show();
                    return;
                }
                $('.state_pass').hide();
                pages = data.pageInfo.TotalPage;
                $('#black_lists').html(temp).find('li').removeBlack();
            });
            if (pages == 1) return;
            $('.page').pagination({
                pages: pages,
                onPageClick: function (pageNumber, event) {
                    $.blackHandler(pageNumber);
                }
            });
        },
        pbY: function (elem, y) {
            var p = elem.css('background-position');
            return p.replace(/\S+$/, y);
        },
        scrollAuto: function () {
            $(window).scroll(function () {
                if ($(window).scrollTop() <= 540) {
                    $('.ssbtn a').attr('href', function () {
                        var path = $(this).attr('href').split('#')[0];
                        return path + '#' + $(window).scrollTop();
                    })
                }
            })
        }
    });
    $.fn.getValid = function (ops) {
        var s = $.extend({}, { url: '/HomeAsync/SendMobile', resend_url: '/HomeAsync/Resend', type: 1 }, ops);
        return this.each(function (i, v) {
            var a = $('#account');
            $(v).bind('click', function () {
                if (a && a.next().hasClass('Validform_right')) {
                    var sec = 60, timer;
                    s.type = ($('.userType') != '') && $('.userType');
                    if (!$(v).hasClass('disable')) {
                        $(v).removeClass('default_btn').addClass('disable');
                        timer = setInterval(function () {
                            $(v).text(--sec + '秒后重新获取');
                            if (sec == 0) {
                                clearInterval(timer);
                                $(v).addClass('default_btn').removeClass('disable');
                                $(v).text('点击获取验证码');
                            }
                        }, 1000);
                        if ($(v).data('send') == 'yes') {
                            $.post(s.resend_url, { type: s.type }, function () {
                                if ($.parseJSON(data).Status == 0) {
                                    $(v).parent().next().show();
                                    $.showmsg('验证码已发送！');
                                    return false;
                                }
                                else {
                                    $.showmsg($.parseJSON(data).Msg);
                                }
                            });
                        }
                        else {
                            $(v).data('send', 'yes');
                            $.post(s.url, { account: $.trim($('#account').val()) }, function (data) {
                                if ($.parseJSON(data).Status == 0) {
                                    $(v).parent().next().show();
                                    $.showmsg('验证码已发送！');
                                    $('.userID').val($.parseJSON(data).Data.id);
                                    $('.userType').val($.parseJSON(data).Data.type);
                                    return false;
                                }
                                else {
                                    $.showmsg($.parseJSON(data).Msg);
                                }
                            });
                        }
                    }
                }
                else if (a && a.val() == '') {
                    $.showmsg('请输入账号');
                    a.focus();
                }
            });
        });
    }
    $.fn.createSelect = function () {
        return this.each(function () {
            var $c = $(this);
            var tag_select = $('<div class="select_box"></div>');
            tag_select.insertBefore($c);

            var select_showbox = $('<div class="select_showbox" style="cursor:pointer"></div>');
            select_showbox.appendTo(tag_select);

            var ul_option = $('<ul class="select_option"></ul>').appendTo(tag_select);
            createOptions(ul_option);

            tag_select.hover(function () {
                ul_option.show();
                ul_option.parent().find(".select_showbox").addClass("active");
            }, function () {
                ul_option.hide();
                ul_option.parent().find(".select_showbox").removeClass("active");
            });

            var li_option = ul_option.find('li');
            li_option.on('click', function () {
                var n = $(this).index(), surl = $c.find('option').eq(n).val();
                $(this).addClass('selected').siblings().removeClass('selected');
                var value = $(this).text();
                select_showbox.text(value);
                switch (n) {
                    case 1:
                        $c.parent().attr('action', surl);
                        break;
                    case 2:
                        $c.parent().attr('action', surl);
                        break;
                    default:
                        $c.parent().attr('action', surl);
                }
                ul_option.hide();
            });

            li_option.hover(function () {
                $(this).addClass('hover').siblings().removeClass('hover');
            }, function () {
                li_option.removeClass('hover');
            });
            function createOptions(ul_list) {
                var options = $c.find('option'),
                    selected_option = options.filter(':selected'),
                    selected_index = selected_option.index(),
                    showbox = ul_list.prev();
                showbox.text(selected_option.text());

                for (var n = 0; n < options.length; n++) {
                    var tag_option = $('<li></li>'),
                        txt_option = options.eq(n).text();
                    tag_option.text(txt_option).appendTo(ul_list);

                    if (n == selected_index) {
                        tag_option.attr('class', 'selected');
                    }
                }
            }
        });
    }
    $.fn.light = function () {
        return this.each(function (i, v) {
            $(v).css({ height: $(document).height(), background: '#000' }).animate({ opacity: 0.7 }, 500);
        });
    }
    $.fn.cm = function () {
        return this.each(function (i, v) {
            $(v).css({ position: 'absolute', zIndex: 999, top: $(window).scrollTop() + Math.round(($(window).height() - $(v).outerHeight()) / 2), left: Math.round(($(window).width() - $(v).outerWidth()) / 2) });
        });
    }
    $.fn.loading = function () {
        return this.each(function () {
            t = $(this);
            var l = $('<div class="loading"></div>').appendTo('body');
            l.find('.loading').css({ width: t.innerWidth(), height: t.innerHeight(), left: t.offset().left, top: offset().top });
        });
    };
    $.fn.valid = function () {
        return this.each(function (i, v) {
            $(v).find('input').each(function () {
                var $c = $(this);
                if ($.trim($c.val()) == '') {
                    $c.siblings('.Validform_title').text('不能为空');
                } else if ($c.attr('data-type') !== undefined) {
                    $.check($.trim($c.val()), $c.attr('data-type')) && $c.text('') || $c.text('格式不正确');
                }
            });
        });
    }
    $.fn.faceOn = function () {
        return this.each(function (i, v) {
            $(v).hover(function () {
                $('<div class="user_pop f_list" id="works_own" action-data="">'
                    + '<div class="user_bg"><img class="cover_bg" src="/Head/20150615/1782662072cc4dfb97a5b21f39f01e98.jpeg" />'
                    + '<div class="user_tx"><a href="#" title="" style="display:inline-block;"><img class="tx" src="/Head/20150618/29ec0f8a4c2f413dbfc898ea953c7531.jpeg" title="" alt="" /></a></div>'
                    + '</div><div class="nick_name"><a href="" title=""><span>中国设计网</span></a></div>'
                    + '<div class="user_attr"><span>平面设计师 河南，郑州</span></div>'
                    + '<div class="user_op">'
                    + '<a href="javascript:void(0);" class="cn_btn add_fol">'
                    + '<i class="unfol"></i>关注</a></div></div>').appendTo('body');
            }, function () {
                $('.user_pop').remove();
            });
        });
    }
    $.fn.modHandler = function (ops) {
        var s = $.extend({ hide: 'yes' }, ops);
        var pic = $('#photo');
        return this.each(function (i, v) {
            $(v).bind('click', function () {
                if (pic.attr('src').length != 0) {
                    subHandler(v);
                } else {
                    $.showmsg("你还没有上传图片");
                }
            });
        });
        function subHandler(v) {
            $(v).addClass('loading');
            $.ajax({
                url: $.picUrl,
                type: 'POST',
                data: { x: pic.data('x'), y: pic.data('y'), w: pic.data('w'), h: pic.data('h'), path: pic.attr('src'), t: 0 },
                success: function (data) {
                    if (s.hide == 'yes') {
                        s.elem.remove();
                    }
                    else {
                        s.elem.addClass('uped').removeClass('edit_tx');
                    }
                    var Data = $.parseJSON(data);
                    if (Data.Status !== 0)
                        return false;
                    $('.cur_tx img').attr('src', Data.Data.LargeHead);
                    $.showmsg('修改头像成功！');
                },
                error: function () {
                    $.showmsg("服务器错误！");
                }
            });
            $(v).removeClass('loading');
        }
    };
    $.fn.pwdS = function () {
        return this.bind('keyup', function () {
            var v = $.trim($(this).val());
            var q = 0;
            if (v.length < 6) {
                q++;
            }
            else if (v.length <= 20) {
                if (/[a-z]+/.test(v)) q++;
                if (/[A-Z]+/.test(v)) q++;
                if (/[0-9]+/.test(v)) q++;
                if (/[^a-zA-Z0-9]+/.test(v)) q++;
                if (q >= 3) q = 3;
            }
            else {
                q = 0;
            }
            $('.pwd_bar').animate({ width: (q * 74) }, 300);
        })
    };
    $.fn.showBind = function () {
        var _t = this;
        return this.each(function (i) {
            var t = window.location.hash.search('#tab') != -1 ? window.location.hash.match(/#tab(\d+)/)[1] : '1';
            if (t) {
                $(_t).eq(t - 1).addClass('active').parent().siblings().find('a').removeClass('active');
                $('.dn_con').eq(t - 1).show().siblings().hide();
            }
            $(this).bind('click', function () {
                $(this).addClass('active').parent().siblings().find('a').removeClass('active');
                $('.dn_con').eq(i).show().siblings().hide();
            })
        })
    }
    $.fn.city = function () {
        return this.each(function (i, v) {
            if ($('.cm').length != 0) return false;
            var f_p = $($.lightbox({ title: '更改所在地', content: '<div class="mb20"><label class="form_m">现在所在地</label></div><div class="city"><select id="province" name="province"></select><label>省</label><select id="city" name="city"></select><label>市</label></div><span class="Validform_checktip"></span></div>', btn: '更改' })).appendTo('body');
            f_p.find('.f_p').cm();
            f_p.find('.sub').bind('click', function () {
                $.post('/HomeAsync/SetAddr', { province: f_p.find('select:eq(0)').val(), city: f_p.find('select:eq(1)').val() }, function (data) {
                    if ($.parseJSON(data).Status == 0) {
                        f_p.remove();
                        $.showmsg($.parseJSON(data).Msg);
                        $('.mod_p .mod_p_txt').text(f_p.find('select:eq(0)').find('option:selected').text() + ',' + f_p.find('select:eq(1)').find('option:selected').text())
                    }
                    else {
                        $.showmsg($.parseJSON(data).Msg);
                    }
                })
            });
            f_p.printCity();
            $('.c').bind('click', function () {
                f_p.remove();
            });
        });
    };
    $.fn.printCity = function () {
        return this.each(function () {
            var pID = $("#province").attr("data-value");
            var cID = $("#city").attr("data-value");
            var p = $.parseJSON($('#cityData').html());
            $(this).find('#province').bind('change', function () { s(); });
            var s1 = $('select')[0];
            var s2 = $('select')[1];
            function s() {
                var loca2 = s1.selectedIndex;
                loca3 = (p[loca2][1]).split("|");
                s2.options.length = 0;
                for (j = 0; j < loca3.length; j++) {
                    s2.options[j] = new Option(loca3[j], j);
                }
            }
            for (i = 0; i < p.length; i++) {
                if (i == pID)
                    s1.options[i] = new Option(p[i][0], i, false, true);
                else
                    s1.options[i] = new Option(p[i][0], i, false, false);
            }
            loca3 = (p[pID][1]).split("|");
            for (l = 0; l < loca3.length; l++) {
                if (l == cID)
                    s2.options[l] = new Option(loca3[l], l, false, true);
                else
                    s2.options[l] = new Option(loca3[l], l, false, false);
            }
        })
    }
    $.fn.showChild = function () {
        return this.each(function () {
            var _t = $(this);
            var _n = _t.find('.per_list_child');
            var timer = null;
            _t.hover(function () {
                _n.show();
            }, function () {
                timer = setTimeout(function () { _n.hide(); }, 200);
            });
            _n.hover(function () {
                clearTimeout(timer);
            }, function () {
                _n.hide();
            })
        })
    };
    //删除作品
    $.fn.deleteWorks = function () {
        return this.each(function (i, v) {
            var temp = $(v).find('a.wol_del').attr('data-href').split('?id=')[1];
            var s = {};
            for (var i = 0; i < temp.length; i++) {
                s[temp[i].split('=')[0]] = temp[i].split('=')[1];
            }
            var f_p = $($.lightbox({ title: '删除作品', content: '<span>你确定要删除这个作品么？</span>', btn: '确定' }));
            $(v).find('a.wol_del').bind('click', function () {
                f_p.appendTo('body').find('.lightBox').light().next('.f_p').cm();
            })
            f_p.find('.c').click(function () { f_p.detach(); });
            f_p.find('.sub').click(function () {
                $.get('/HomeWorks/Remove?id=' + temp, function (data) {
                    window.location.reload();
                });
            })
        })
    }
    //分享作品到大圣创易
    $.fn.shareWorks = function () {
        var me = this;
        var worksClass = {};
        // 获取作品分类
        $.get('/Api/Dsook/WorksClass', function (data) {
            if(data.Status == 0){
                worksClass = data.Data;
                var optStr = '';
                for (var key in worksClass) {
                    optStr += '<option value="' + worksClass[key].Key + '">' + worksClass[key].Value + '</option>'
                }
                return me.each(function (i, v) {
                    var temp = $(v).find('a.wol_del').attr('data-href').split('?id=')[1];
                    var s = {};
                    for (var i = 0; i < temp.length; i++) {
                        s[temp[i].split('=')[0]] = temp[i].split('=')[1];
                    }
                    var f_p = $($.lightbox({
                        title: '分发作品至大圣创易平台',
                        content: '<span>分发分类：</span>'+
                        '<select id="shareClass_">'+ optStr + '</select>',
                        btn: '确定' }));
                    var t_btn = $(v).find('.share_btn');
                    $(v).find('.share_btn').bind('click', function () {
                        f_p.appendTo('body').find('.lightBox').light().next('.f_p').cm();
                    });
                    f_p.find('.c').click(function () {
                        f_p.detach();
                    });
                    f_p.find('.sub').click(function () {
                        $.get('/HomeWorksAsync/Share?worksID='+temp+'&shareClass='+$('#shareClass_').val(), function (data) {
                            if(data.Status == 0){
                                f_p.detach();
                                $.showmsg('分发成功！');
                                t_btn.addClass('hide');
                                t_btn.next('.shared').removeClass('hide');
                            }else{
                                $.showmsg(data.Msg);
                            }
                        });
                    });
                });
            }
        });
    }
    //取消关注
    $.fn.unfollow = function () {
        return this.each(function (i, v) {
            var temp = $(v).attr('action-data').split('&');
            var s = {};
            for (var i = 0; i < temp.length; i++) {
                s[temp[i].split('=')[0]] = temp[i].split('=')[1];
            }
            temp = [];
            var f_p = $($.lightbox({ title: '取消关注', content: '<span>你确定取消对<b style="color:green">' + s.nickName + '</b>的关注么？</span>', btn: '确定' }));
            $(v).find('a.cancel_fol').bind('click', function () {
                if ($(this).attr('do-type') == 1) {
                    f_p.appendTo('body').find('.lightBox').light().next('.f_p').cm();
                }
                else {
                    $.post('/HomeAsync/CancelFollow', { followID: s.userID }, function (data) {
                        if ($.parseJSON(data).Status === 0) {
                            $(v).removeClass('ed');
                            $('#f_count').text(Number($('#f_count').text()) - 1);
                            return false;
                        }
                        else {
                            $.showmsg($.parseJSON(data).Msg);
                        }
                    });
                }
            })
            f_p.find('.c').click(function () { f_p.detach(); });
            f_p.find('.sub').click(function () {
                $.post('/HomeAsync/CancelFollow', { followID: s.userID }, function (data) {
                    if ($.parseJSON(data).Status === 0) {
                        $('#f_count').text(Number($('#f_count').text()) - 1);
                        f_p.remove();
                        $(v).remove();
                        //$.showmsg('你已经取消对' + s.nickName + '的关注');
                        window.location.reload();
                    }
                    else { $.showmsg($.parseJSON(data).Msg); }
                });
            });
        });
    }
    //移除粉丝
    $.fn.unfans = function () {
        return this.each(function () {
            _t = $(this);
            var temp = _t.attr('action-data').split('&');
            var s = {};
            for (var i = 0; i < temp.length; i++) {
                s[temp[i].split('=')[0]] = temp[i].split('=')[1];
            }
            temp = [];
            var op = $('<div></div>');
            var f_p = $($.lightbox({ title: '移除粉丝', content: '<span>你确定要移除<b style="color:green">' + s.nickName + '</b>？</span>', btn: '确定' }));
            _t.find('a[action-type=remove_fans]').bind('click', function () {
                f_p.appendTo('body').find('.lightBox').light().next('.f_p').cm();
            })
            f_p.find('.c').click(function () { f_p.detach(); });
            f_p.find('.sub').click(function () {
                $.post('/HomeAsync/RemoveFans', { fansID: s.userID }, function (data) {
                    if ($.parseJSON(data).Status === 0) {
                        $('#f_count').text(Number($('#f_count').text()) - 1);
                        f_p.detach();
                        _t.remove();
                        location.reload(true);
                    }
                    else { $.showmsg($.parseJSON(data).Msg); }
                });
            });
        });
    }
    //添加关注
    $.fn.addFollow = function () {
        return this.each(function (i, v) {
            var temp = $(v).attr('action-data').split('&');
            var s = {}, _t, h;
            for (var i = 0; i < temp.length; i++) {
                s[temp[i].split('=')[0]] = temp[i].split('=')[1];
            }
            temp = [];
            $(v).find('.add_fol').bind('click', function () {
                _t = $(this)
                h = _t.html();
                _t.addClass('load').html('<i class="unfol"></i>关注中');
                $.post('/HomeAsync/Follow', { followID: s.userID }, function (data, textStatus, xhr) {
                    var res = $.parseJSON(xhr.getResponseHeader('X-Responded-JSON'));
                    if (res && res.status === 401) {
                        var s = res.headers.location.toLocaleLowerCase();
                        window.location.href = s.substring(0, s.search('=')) + '=' + window.location.pathname + window.location.search;
                        _t.html(h).removeClass('load');
                        return;
                    }
                    else {
                        var d = $.parseJSON(data);
                        if (d.Status === 0) {
                            if ($(v).find('.add_fol').html('<div class="iUnfol"><i class="unfol">+</i>关注</div>').attr('do-type') == 1) {
                                $(v).addClass('ed');
                            }
                            else {
                                _t.removeClass('add_fol load').addClass('def').unbind();
                                // d.Data.type == '1' ? _t.html('<div class="iFoled"><i class="foled"></i>已关注</div>') : _t.html('<div class="iFolEach"><i class="fol_each"></i>互相关注</div>');
                                d.Data.type == '1' ? _t.html('<div class="iFoled"><i class="foled"></i>已关注</div>') : _t.html('<div class="iFolEach">互相关注</div>');
                            }
                            $('#f_count').text(Number($('#f_count').text()) + 1);
                            try {
                                msg.server.addOtherNotice(d.Data.count, d.Data.userID, 1);
                            }
                            catch (e) { }
                        }
                        else {
                            _t.html(h).removeClass('load');
                            $.showmsg(d.Msg);
                        }
                    }
                });
            });
        });
    }
    //加入黑名单
    $.fn.addBlack = function () {
        return this.each(function (i, v) {
            var temp = $(v).attr('action-data').split('&');
            var s = {};
            for (var i = 0; i < temp.length; i++) {
                s[temp[i].split('=')[0]] = temp[i].split('=')[1];
            }
            temp = [];
            var op = $('<div></div>');
            var f_p = $($.lightbox({ title: '加入黑名单', content: '<span>确认将<b style="color:green">' + s.nickName + '</b>加入到黑名单中么？</span><p>你和他将自动解除关注关系, 将不能继续关注！</p>', btn: '加入黑名单' }));
            $(v).find('a[action-type=add_black]').bind('click', function () {
                f_p.appendTo('body').find('.lightBox').light().next('.f_p').cm();
            });
            f_p.find('.c').bind('click', function () { f_p.detach(); });
            f_p.find('.sub').click(function () {
                f_p.detach();
                $.post('/HomeAsync/ToBlack', { ID: s.userID }, function (data, textStatus, xhr) {
                    var d = $.parseJSON(data);
                    if (d.Status === 0) {
                        $(v).removeClass('ed');
                        $.showmsg('您已经将' + s.nickName + '加入到黑名单中！');
                    }
                    else { $.showmsg(d.Msg); }
                });
            });
        });
    };
    //移除黑名单
    $.fn.removeBlack = function () {
        return this.each(function (i, v) {
            var temp = $(v).attr('action-data').split('&');
            var s = {};
            for (var i = 0; i < temp.length; i++) {
                s[temp[i].split('=')[0]] = temp[i].split('=')[1];
            }
            temp = [];
            $(v).find('.s_btn').bind('click', function () {
                var f_p = $($.lightbox({ title: '移除黑名单', content: '<span>确认将<b style="color:green">' + s.nickName + '</b>移除黑名单？</span>', btn: '移除黑名单' }));
                f_p.appendTo('body').find('.lightBox').light().next('.f_p').cm();
                f_p.find('.c').bind('click', function () { f_p.detach(); });
                f_p.find('.sub').click(function () {
                    $.post('/HomeAsync/RemoveBlack', { ID: s.userID }, function (data) {
                        if ($.parseJSON(data).Status === 0) {
                            f_p.remove();
                            $(v).remove();
                            $.showmsg('您已经将' + s.nickName + '移除黑名单！');
                            $.blackHandler($(v).children().length == 0 ? $('#black_lists').data('page') - 1 : $('#black_lists').data('page'));
                        }
                    });
                });
            });
        });
    }
    $.fn.random = function () {
        return this.each(function (i, v) {
            var at = parseInt($(v).attr('action-type'));
            var temp = '<li class="p bt clear" action-data="userID=$UserID$&nickName=$NickName$">'
                + '<a href="$Url$" target="_blank" class="fl"><img src="$UserHead$"></a>' +
                '<div class="clear"><span><a href="$Url$" target="_blank" class="s_C nick_name">$NickName$</a></span><em>$skill$</em>' +
                '<div class="br">' +
                '<a class="follow add_fol" href="javascript:void(0);"><i class="fol_i"></i>关注</a>' +
                '</div></div>' +
                '</li>'
            $.get('/HomeAsync/GetInterestingUser', { type: at, page: $(v).data('page') }, function (data) {
                var d = $.parseJSON(data);
                if (d.Data.length === 0) {
                    $.showmsg('没有更多用户了！');
                    $(v).addClass('disabled');
                }
                else {
                    $(v).data('page', parseInt($(v).data('page')) + 1).parent().next().html(function () {
                        var str = '';
                        for (var i = 0; i < d.Data.length; i++) {
                            str += temp.replace(/\$([A-Za-z]+)?\$/g, function (a) {
                                return d.Data[i][a.substring(1, a.length - 1)];
                            });
                        }
                        return str;
                    }).children().addFollow();
                }
            });
        });
    }
    $.fn.moveUp = function () {
        return this.each(function (i, v) {
            var _t = $(v);
            var _n = _t.find('.slideUp');
            var _o = { top: _n.attr('data-top') + 'px', down: _n.attr('data-down') + 'px' };
            _t.hover(function () {
                _n.stop().animate({ top: _o.top }, 200);
            }, function () {
                _n.stop().animate({ top: _o.down }, 200);
            });
        });
    }


})(jQuery, window, document);
$(document).ready(function () {
    /*公共*/
    $('#fol_list>li').unfollow();
    $('#fans_list>li, #coms_list>li').unfans().addFollow();
    $('#other_op').addFollow().unfollow().addBlack();
    $('.back').css('right', ($(window).width() - 1200) / 2 - 70 + 'px').bind('click', function () {
        $('html,body').animate({ 'scrollTop': 0 }, 500, function () {
            $('.back').fadeOut();
        });
        $(this).find('.help').bind('click', function (e) {
            e.stopPropagation();
        });
    });
    $(window).bind({
        'scroll': function () {
            var st = $(window).scrollTop();
            var t = $('.dingwei-top').length && $('.dingwei-top').offset().top;
            var b = $('.dingwei-bottom').length && $('.dingwei-bottom').offset().top;
            if (st > 100) {
                $('.back').fadeIn();
            }
            else {
                $('.back').fadeOut();
            }
            if (st <= (b - $('#erwei').height()) && st > t) {
                $('#erwei').css({ position: 'fixed', top: 80, bottom: '0px' });
            } else {
                $('#erwei').css({ position: 'relative', top: '0px', bottom: 'auto' });
            }
        }
    });

    $('#dn_navs a').showBind();

    typeof $().lazyload === 'function' && $('.crop_img img').lazyload({
        threshold: 200,
        effect: "show",
        tolerance: 60,
        failure_limit: 10,
        container: window,
        load: function () {
            $(this).addClass('white_bg');
            if($(this).attr('data-original') == $(this).attr('src')) {
                $(this).parent('.crop_img').css('background', 'none');
            }
        },
        placeholder: PicHub + '/Content/images/blank.gif'
    });
    $('.view-list').hover(function () {
        $('.hasCount') && $('.hasCount').hide();
        $(this).addClass('hover').find('.s-C').show();
    }, function () {
        $('.hasCount').show();
        $(this).removeClass('hover').find('.s-C').hide();
    });
    $('.fn_nav_child').hover(function () {
        $(this).find('.f_icon').stop().animate({ 'backgroundPositionY': '-61px', 'backgroundColor': '#333024' }, 300);
    }, function () {
        $(this).find('.f_icon').stop().animate({ 'backgroundPositionY': '0px', 'background-color': '#fff' }, 300);
    });
    $('.more').showChild();
    $('#sites').children().each(function (i, v) {
        $(v).hover(function () {
            $('#sites').next().children().eq(i).animate({ top: 60, opacity: 1 }).show();
        }, function () {
            $('#sites').next().children().eq(i).animate({ top: 40, opacity: 0 }).hide();
        });
    });
    $('#choose').createSelect();
    checkCookie();

    //个人中心-资料下方的导航栏 点击加载页面后自动滚动到该位置
    $.scrollAuto();     //只要发生滚动就给a的链接加上参数
    if (location.hash == '') {
        $('.ssbtn a').attr('href', function () {
            var path = $(this).attr('href').split('#')[0];
            return path;
        })
    } else if (location.hash != '' && location.hash != '#540') {
        $(window).scrollTop(location.hash.split('#')[1]);
        // if (location.hash.split('#')[1] != "1" && location.hash.split('#')[1] != "2" && location.hash.split('#')[1] != "3" && location.hash.split('#')[1] != "4" && location.hash.split('#')[1] != "5") {
        if (location.hash.split('#')[1] >=10) {
            $('html,body').animate({ scrollTop: '540px' }, 200);
        }
    } else if (location.hash == '#540') {
        $(window).scrollTop(540);
    }
});


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
    if (!window.location.origin) {
        window.location.origin = window.location.protocol + "//" + window.location.host;
    }
    if($(window).scrollTop() <= 120){
        $('.nav_logo').attr('href', window.location.origin);
    }
});

$(function () {
    //鼠标移入图片放大
    if ($('.zoom').zoomImgRollover != undefined) {
        $('.zoom').zoomImgRollover();
    }

    $('.paging_a_selected').on('click', function (e) { //阻止当前页点击跳转
        e.preventDefault();
    });
    //    分页跳转
    $('#jump_go').on('click', function () {
        var href = $('.paging_a_selected').attr('href'),
            currentPage = +$('.paging_a_selected').html(),
            maxPage = +$('.triangle_right').parents('.paging_lists').prev().find('a').html(),
            newPage = +$('#jump_page').val(),
            idx = href.lastIndexOf('=');
        href = href.slice(0, idx + 1);
        newPage = parseInt(newPage);
        if (typeof newPage == "number" && newPage >= 1 && newPage <= maxPage) {
            // window.location.href = href + newPage;
            $('#pageForm').submit();
        } else {
            $('#jumpGO_isMax').show();
            setTimeout(function () {
                $('#jumpGO_isMax').hide();
            }, 2000);
        }
    })
})

function getCookie(c_name) {
    if (document.cookie.length > 0) {
        c_start = document.cookie.indexOf(c_name + "=")
        if (c_start != -1) {
            c_start = c_start + c_name.length + 1
            c_end = document.cookie.indexOf(";", c_start)
            if (c_end == -1) c_end = document.cookie.length
            return unescape(document.cookie.substring(c_start, c_end))
        }
    }
    return ""
}

function setCookie(c_name, value, expiredays) {
    var exdate = new Date()
    exdate.setDate(exdate.getDate() + expiredays)
    document.cookie = c_name + "=" + escape(value) +
        ((expiredays == null) ? "" : "; expires=" + exdate.toGMTString())
}

function checkCookie() {
    if (getCookie('notice') == 'hide') {
        return;
    } else {
        $('.web-notice').show();
        $('.web-notice .hide').bind('click', function () {
            $('.web-notice').hide();
            setCookie('notice', 'hide', 30);
        })
    }
}
