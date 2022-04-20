// $(function() {

//     var websocket = null;

//     var webSocket = createSocketConnection("ws://im.nihongdengxia.com/im","eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodHRwOlwvXC9pbS5uaWhvbmdkZW5neGlhLmNvbVwvIiwiYXVkIjoiaHR0cDpcL1wvaW0ubmlob25nZGVuZ3hpYS5jb21cLyIsImlhdCI6MTYyMzAzNzE4NCwidXNlciI6MX0.7s0gxXlal1Hqbz9naKmlo2fgfNNmz_73jcfySSE4s-k");
//     socketEvent(webSocket);

//     function send() {
//         output(res, 'toMessage');
//         // var s = {
//         //     mine: i,
//         //     to: a.data
//         // }
//         // var o = {
//         //     username: s.mine.username,
//         //     avatar: s.mine.avatar,
//         //     id: s.to.id,
//         //     type: s.to.type,
//         //     content: s.mine.content,
//         //     timestamp: (new Date).getTime(),
//         //     mine: !0
//         // };

//         var o = {
//             avatar: "https://s.gravatar.com/avatar/b04a7799bfc108fea09eef6270505090",
//             content: "123",
//             id: 1,
//             mine: true,
//             username: "chenli@qq.com",
//         }
//         var to = {
//             "avatar":"https://s.gravatar.com/avatar/487f87505f619bf9ea08f26bb34f8118",
//             "id": 2,
//             "name": "123@qq.com",
//             "sign": "",
//             "status": "offline",
//             "type": "friend",
//             "username": "123@qq.com",
//             "mine":o
//         }
//         let cmd = (to.type === 'friend') ? friend_send_cmd : group_send_cmd;
//         let data = {
//           message_id: messageId(),
//           from_user_id: to.mine.id,
//           to_id: parseInt(to.id),
//           content: to.mine.content
//         };
//         let msg = createMessage(cmd, data);
//         wsSend(msg);
// 　　 //浏览器是否支持websocket
//     // if ("WebSocket" in window) {
//     //     try {
//     //         //websocket = new WebSocket($("#socketUrl").val());
//     //         websocket = new WebSocket("ws://im.nihongdengxia.com/im","eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodHRwOlwvXC9pbS5uaWhvbmdkZW5neGlhLmNvbVwvIiwiYXVkIjoiaHR0cDpcL1wvaW0ubmlob25nZGVuZ3hpYS5jb21cLyIsImlhdCI6MTYyMzAzNzE4NCwidXNlciI6MX0.7s0gxXlal1Hqbz9naKmlo2fgfNNmz_73jcfySSE4s-k");
//     //         websocket.onopen = function(data) {
//     //             console.log("connect success :" + data.data);
//     //         };

//     //         websocket.onmessage = function(data) {
//     //             console.log("onmessage : " + data.data);
//     //         };

//     //         websocket.onerror = function(e) {
//     //             console.log('connect error :'+e);
//     //         };

//     //         websocket.onclose = function(data) {
//     //             console.log('connect error :'+data);
//     //         };

//     //         window.onbeforeunload = function() {
//     //             websocket.close();
//     //         }

//     //     } catch (error) {
//     //         console.log('connect error :'+error);
//     //     }

//     // } else {
//     //     console.log("Websockets not supported");
//     // }

//     // $('.close').click(function() {
//     //     if (websocket != null) {
//     //         websocket.close();
//     //     }
//     // });

// });

import {createSocketConnection,createMessage,wsSend,socketEvent,askMsg} from "./im.js";
import {output} from "./util.js";
import {
    friend_send_cmd,
    group_send_cmd,
    user_ping_cmd,
    system_error_cmd,
    system_event_cmd,
    friend_get_unread_message_cmd,
    user_get_unread_application_count_cmd
  } from "./api.js";
$(function() {

    var webSocket = createSocketConnection("ws://im.nihongdengxia.com/im","eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodHRwOlwvXC9pbS5uaWhvbmdkZW5neGlhLmNvbVwvIiwiYXVkIjoiaHR0cDpcL1wvaW0ubmlob25nZGVuZ3hpYS5jb21cLyIsImlhdCI6MTYyMzAzNzE4NCwidXNlciI6MX0.7s0gxXlal1Hqbz9naKmlo2fgfNNmz_73jcfySSE4s-k");
    socketEvent(webSocket);


    $('#ack').click(send)
    $(".msg").keydown(function(e){
        if(e.keyCode == 13){
            return send();
        }
    });
})

// function reply(headSrc, str) {
//     var html = `<object class="left">
//             <a href="javascript:;">
//                 <div class="aui-flex fix-header-style detail-sub-comtent">
//                     <div class="topic-user-img">
//                         <object>
//                             <a href="404.html">
//                                 <img src="images/pd-001.png" alt>
//                             </a>
//                         </object>
//                     </div>
//                     <div class="topic-user-box">
//                         <!-- <h3>
//                             <object>
//                                 <a href="404.html">职业第三者</a>
//                             </object>
//                         </h3> -->
//                         <div class="detail-sub-comment-content">
//                             <p>杨庄路是一条典型的通州街道，门脸窄小的SPA店.</p>
//                         </div>
//                     </div>
//                 </div>
                
//             </a>
//         </object>`;
//     return upView(html);
// }
// function ask(headSrc, str) {
//     var html = `
//         <object class="right">
//             <a href="javascript:;">
//                 <div class="aui-flex fix-header-style detail-sub-comtent">
//                     <div class="topic-user-img">
//                         <object>
//                             <a href="404.html">
//                                 <img src="images/pd-001.png" alt>
//                             </a>
//                         </object>
//                     </div>
//                     <div class="topic-user-box">
//                         <h3>
//                             <object>
//                                 <a href="404.html">职业第三者</a>
//                             </object>
//                         </h3>
//                         <div class="detail-sub-comment-content">
//                             <p>`+str+`</p>
//                         </div>
//                     </div>
//                 </div>
                
//             </a>
//         </object>`;
//     return upView(html);
// }
// function upView(html) {
//     let message = $('.chat-content');
//     message.append(html);
//     // console.log(message.outerHeight() - window.innerHeight)
//     return $('.aui-scrollView').animate({
//         scrollTop: message.outerHeight() - window.innerHeight + 15
//     }, 200);
// }

// function send(msg){
//     if(canSend){
//         let input = $(".msg");
//         ask("", input.val());
//         input.val('');
//         test();
//     }
//     $(".msg")[0].focus();
// }
function messageId() {
    return Number(Math.random().toString().substr(3, 3) + Date.now()).toString(36)
  }
function send() {

    var o = {
        avatar: "https://s.gravatar.com/avatar/b04a7799bfc108fea09eef6270505090",
        content:  $(".msg").val(),
        id: 1,
        mine: true,
        username: "职业第三者",
    }
    var to = {
        "avatar":"https://s.gravatar.com/avatar/487f87505f619bf9ea08f26bb34f8118",
        "id": 2,
        "name": "123@qq.com",
        "sign": "",
        "status": "offline",
        "type": "friend",
        "username": "123@qq.com",
        "mine":o
    }
    let cmd = (to.type === 'friend') ? friend_send_cmd : group_send_cmd;
    let data = {
    message_id: messageId(),
    from_user_id: to.mine.id,
    to_id: parseInt(to.id),
    content: to.mine.content
    };
    let msg = createMessage(cmd, data);
    wsSend(msg);
    askMsg(o)
    $(".msg").val('');
    $(".msg")[0].focus();
}

