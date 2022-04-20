import {output, getCookie} from "./util.js";

function getRequest(url, params, success_callback, fail_callback) {
  return Request('get', url, params, success_callback, fail_callback)
}

function postRequest(url, params, success_callback, fail_callback) {
  return Request('post', url, params, success_callback, fail_callback)
}


function Request(type, url, params, success_callback, fail_callback) {
    $.ajax({
      url: url,
      data: params,
      type: type,
      beforeSend: function (request) {
        request.setRequestHeader("Authorization", 'Bearer ' + getCookie("IM_TOKEN"));
      },
      success: function (data) {
        output(data, url);
        //layer.closeAll('loading');
        if (JSON.stringify(data) === '{}'){
          return false;
        }
        if (data.code && data.code != 0) {
          //layer.msg(data.code + ' : ' + data.msg);
          alert(data.code + ' : ' + data.msg)
          fail_callback && fail_callback(data.data, data.msg);
          return false;
        }
        alert(data.msg, {time: 500});
        success_callback && success_callback(data.data, data.msg);
      },
      error: function () {
        //layer.closeAll('loading');
        alert('Interface cannot connect : ' + url)
      }
    })
}

export {
  getRequest,
  postRequest
}
