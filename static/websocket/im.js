import {output, isEmpty, getCookie} from "./util.js";
import {
  user_ping_cmd,
  system_error_cmd,
  system_event_cmd,
  friend_get_unread_message_cmd,
  user_get_unread_application_count_cmd,
  friend_read_msg_cmd
} from "./api.js";

var MessageActive = {
  setUserStatus: function (data) {
    output(data, 'setUserStatus');
    //layui.layim.setFriendStatus(data.user_id, data.status)
  },
  getMessage: function (data) {
    output(data, 'getMessage');
    //layui.layim.getMessage(data);
    reply(data)


    if (data.type === 'friend') {
      let msg = createMessage(friend_read_msg_cmd, {
        'message_id': data.cid
      });
      wsSend(msg)
    }
  },
  onlineNumber: function (data) {
    //layui.jquery("#onlineNumber").html(data)
  },
  getUnreadApplicationCount: function (data) {
    if (data === 0) return false;
    // layui.layim.msgbox(data)
  },
  friendAgreeApply: function (data) {
     output(data, 'friendAgreeApply');
  },
  groupAgreeApply: function (data) {
    output(data, 'groupAgreeApply');
  },
  friendVideoRoom: function (data) {
    output(data, 'friendVideoRoom');
  }
};

var Im;
var heartbeat;
var messageList = {};

function createSocketConnection(url, protocols) {
  output(url, 'createSocketConnection');
  console.log(url);
  console.log(protocols);
  Im = new WebSocket(url, protocols);
  return Im;
}

function createMessage(cmd, data = {}, ext = {}) {
  let msg = {
    cmd: cmd,
    data: data,
    ext: ext
  };
  output(msg);
  ack(msg);
  return JSON.stringify(msg);
}

function ack(msg, num = 1) {
  let data = msg.data;
  let message_id = data.message_id;
  if (isEmpty(message_id)) return false;
  messageList[message_id] = {
    msg: msg,
    num: num,
    timer: setTimeout(function () {
      output(num, message_id + '的发送次数');
      if (num > 3) {
        if (!isEmpty(data.content)) {
          output(num, '重试次数大于3进行提示');

          console.log('消息发送失败：' + data.content)
          alert('消息发送失败：' + data.content)
          // layui.layer.msg('消息发送失败：' + data.content, {
          //   time: 0
          //   , btn: ['重试', '取消']
          //   , yes: function (index) {
          //     Im.send(JSON.stringify(msg));
          //     ack(messageList[message_id].msg, messageList[message_id]['num'] + 1);
          //     layui.layer.close(index);
          //   },
          //   btn2: function (index) {
          //     delete messageList[message_id];
          //     layui.layer.close(index);
          //   }
          // });
        }
      } else {
        Im.send(JSON.stringify(msg));
        ack(messageList[message_id].msg, messageList[message_id]['num'] + 1);
      }
    }, 5000)
  };
  output(messageList);
}


function wsOpen(event) {
  output(event, 'onOpen');
  heartbeat = setInterval(function () {
    wsSend(createMessage(user_ping_cmd));
  }, 10000);
  infoInit();
}

function infoInit() {
  wsSend(createMessage(friend_get_unread_message_cmd, {}));
  wsSend(createMessage(user_get_unread_application_count_cmd, {}))
}

function wsReceive(event) {
  let result = eval('(' + event.data + ')');
  output(result, 'onMessage');

  if (JSON.stringify(result) === '{}'){
    return false;
  }

  // if (layui.jquery.isEmptyObject(result)) {
  //   return false;
  // }

  if (result.cmd && result.cmd === system_error_cmd) {
    // layer.msg(result.cmd + ' : ' + result.msg);
    console.log(result.cmd + ' : ' + result.msg)
    clearMessageListTimer(result);
    return false;
  }

  if (result.cmd && result.cmd === system_event_cmd) {
    let method = result.method;
    MessageActive[method] ? MessageActive[method](result.data) : '';
    return false;
  }


  if (result.cmd && result.cmd === user_ping_cmd) {
    return false;
  }

  clearMessageListTimer(result);

}

function clearMessageListTimer(result) {
  let message_id = result.data.message_id ? result.data.message_id : '';
  if (message_id === '') return false;
  console.log(message_id)
  clearInterval(messageList[message_id].timer);
  delete messageList[message_id];
}

function wsError(event) {
  output(event, 'onError');
  clearInterval(heartbeat);
  reloadSocket(event);
}

function wsClose(event) {
  output(event, 'onClose');
  clearInterval(heartbeat);
  reloadSocket(event)
}

function reloadSocket(event) {
  console.log('连接异常关闭')
  alert('连接异常关闭')
  // layui.layer.msg(event.reason, {
  //   time: 0
  //   , title: '连接异常关闭'
  //   , btn: ['重试', '取消']
  //   , yes: function (index) {
  //     var wsUrl = layui.jquery(".wsUrl").val();
  //     Im = createSocketConnection(wsUrl, getCookie('IM_TOKEN'));
  //     socketEvent(Im);
  //     layui.layer.close(index);
  //   },
  //   btn2: function (index) {
  //     layui.layer.close(index);
  //   }
  // });
}

function wsSend(data) {
  Im.send(data)
}

function socketEvent(webSocket) {
  webSocket.onopen = function (event) {
    wsOpen(event);
  };
  webSocket.onmessage = function (event) {
    wsReceive(event);
  };
  webSocket.onerror = function (event) {
    wsError(event)
  };
  webSocket.onclose = function (event) {
    wsClose(event)
  };
}


// 接收对方消息
function reply(data) {
  var html = `<object class="left">
          <a href="javascript:;">
              <div class="aui-flex fix-header-style detail-sub-comtent">
                  <div class="topic-user-img">
                      <object>
                          <a href="404.html">
                              <img src="`+data.avatar+`" alt>
                          </a>
                      </object>
                  </div>
                  <div class="topic-user-box">
                      <!-- <h3>
                          <object>
                              <a href="404.html">`+data.username+`</a>
                          </object>
                      </h3> -->
                      <div class="detail-sub-comment-content">
                          <p>`+data.content+`</p>
                      </div>
                  </div>
              </div>
              
          </a>
      </object>`;
  return upView(html);
}
// 自己发送的信息
function askMsg(data) {
  var html = `
      <object class="right">
          <a href="javascript:;">
              <div class="aui-flex fix-header-style detail-sub-comtent">
                  <div class="topic-user-img">
                      <object>
                          <a href="404.html">
                              <img src="`+data.avatar+`" alt>
                          </a>
                      </object>
                  </div>
                  <div class="topic-user-box">
                      <h3>
                          <object>
                              <a href="404.html">`+data.username+`</a>
                          </object>
                      </h3>
                      <div class="detail-sub-comment-content">
                          <p>`+data.content+`</p>
                      </div>
                  </div>
              </div>
              
          </a>
      </object>`;
  return upView(html);
}
// 消息上拉
function upView(html) {
  let message = $('.chat-content');
  message.append(html);
  // console.log(message.outerHeight() - window.innerHeight)
  return $('.aui-scrollView').animate({
      scrollTop: message.outerHeight() - window.innerHeight + 15
  }, 200);
}

export {
  createSocketConnection,
  socketEvent,
  wsOpen,
  wsReceive,
  wsError,
  wsClose,
  wsSend,
  createMessage,
  askMsg
}
