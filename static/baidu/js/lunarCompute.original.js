var lunarInfo = new Array(
    0x04bd8, 0x04ae0, 0x0a570, 0x054d5, 0x0d260, 0x0d950, 0x16554, 0x056a0, 0x09ad0, 0x055d2,
    0x04ae0, 0x0a5b6, 0x0a4d0, 0x0d250, 0x1d255, 0x0b540, 0x0d6a0, 0x0ada2, 0x095b0, 0x14977,
    0x04970, 0x0a4b0, 0x0b4b5, 0x06a50, 0x06d40, 0x1ab54, 0x02b60, 0x09570, 0x052f2, 0x04970,
    0x06566, 0x0d4a0, 0x0ea50, 0x06e95, 0x05ad0, 0x02b60, 0x186e3, 0x092e0, 0x1c8d7, 0x0c950,
    0x0d4a0, 0x1d8a6, 0x0b550, 0x056a0, 0x1a5b4, 0x025d0, 0x092d0, 0x0d2b2, 0x0a950, 0x0b557,
    0x06ca0, 0x0b550, 0x15355, 0x04da0, 0x0a5d0, 0x14573, 0x052d0, 0x0a9a8, 0x0e950, 0x06aa0,
    0x0aea6, 0x0ab50, 0x04b60, 0x0aae4, 0x0a570, 0x05260, 0x0f263, 0x0d950, 0x05b57, 0x056a0,
    0x096d0, 0x04dd5, 0x04ad0, 0x0a4d0, 0x0d4d4, 0x0d250, 0x0d558, 0x0b540, 0x0b5a0, 0x195a6,
    0x095b0, 0x049b0, 0x0a974, 0x0a4b0, 0x0b27a, 0x06a50, 0x06d40, 0x0af46, 0x0ab60, 0x09570,
    0x04af5, 0x04970, 0x064b0, 0x074a3, 0x0ea50, 0x06b58, 0x055c0, 0x0ab60, 0x096d5, 0x092e0,
    0x0c960, 0x0d954, 0x0d4a0, 0x0da50, 0x07552, 0x056a0, 0x0abb7, 0x025d0, 0x092d0, 0x0cab5,
    0x0a950, 0x0b4a0, 0x0baa4, 0x0ad50, 0x055d9, 0x04ba0, 0x0a5b0, 0x15176, 0x052b0, 0x0a930,
    0x07954, 0x06aa0, 0x0ad50, 0x05b52, 0x04b60, 0x0a6e6, 0x0a4e0, 0x0d260, 0x0ea65, 0x0d530,
    0x05aa0, 0x076a3, 0x096d0, 0x04bd7, 0x04ad0, 0x0a4d0, 0x1d0b6, 0x0d250, 0x0d520, 0x0dd45,
    0x0b5a0, 0x056d0, 0x055b2, 0x049b0, 0x0a577, 0x0a4b0, 0x0aa50, 0x1b255, 0x06d20, 0x0ada0)

var solarMonth = new Array(31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31);
var Animals = new Array("鼠", "牛", "虎", "兔", "龙", "蛇", "马", "羊", "猴", "鸡", "狗", "猪");
var solarTerm = new Array("小寒", "大寒", "立春", "雨水", "惊蛰", "春分", "清明", "谷雨", "立夏", "小满", "芒种", "夏至", "小暑", "大暑", "立秋", "处暑", "白露", "秋分", "寒露", "霜降", "立冬", "小雪", "大雪", "冬至");
var sTermInfo = new Array(0, 21208, 42467, 63836, 85337, 107014, 128867, 150921, 173149, 195551, 218072, 240693, 263343, 285989, 308563, 331033, 353350, 375494, 397447, 419210, 440795, 462224, 483532, 504758);
var nStr1 = new Array('日', '一', '二', '三', '四', '五', '六', '七', '八', '九', '十');
var nStr2 = new Array('初', '十', '廿', '卅');


// var Today = new Date();
var localTimezoneOffset = new Date().getTimezoneOffset()
var Today = window.beijingTime;
var tY = Today.getFullYear();
var tM = Today.getMonth();
var tD = Today.getDate();

function lYearDays(y) {
    var i, sum = 348;
    for (i = 0x8000; i > 0x8; i >>= 1) sum += (lunarInfo[y - 1900] & i) ? 1 : 0;
    return (sum + leapDays(y));
}

//返回农历y年闰月的天数
function leapDays(y) {
    if (leapMonth(y)) return ((lunarInfo[y - 1900] & 0x10000) ? 30 : 29);
    else return (0);
}

//判断y年的农历中那个月是闰月,不是闰月返回0
function leapMonth(y) {
    return (lunarInfo[y - 1900] & 0xf);
}

//返回农历y年m月的总天数
function monthDays(y, m) {
    return ((lunarInfo[y - 1900] & (0x10000 >> m)) ? 30 : 29);
}

//算出当前月第一天的农历日期和当前农历日期下一个月农历的第一天日期
function Dianaday(objDate) {
    let f = 1;
    let o = localTimezoneOffset / 60 + 8;
    if (o == -6) {
        f = 2;
    }else if (o < 0){
        f = 1;
    }else if (o == 0){
        f = 0;
    }
    objDate = new Date(objDate + f * 60 * 60 * 20 * 1000).getTime()
    var i, leap = 0, temp = 0;
    var baseDate = new Date(1900, 0, 31);
    var offset = (objDate - baseDate) / 86400000;
    this.dayCyl = offset + 40;
    this.monCyl = 14;
    for (i = 1900; i < 2050 && offset > 0; i++) {
        temp = lYearDays(i)
        offset -= temp;
        this.monCyl += 12;
    }
    if (offset < 0) {
        offset += temp;
        i--;
        this.monCyl -= 12;
    }
    this.year = i;
    this.yearCyl = i - 1864;
    leap = leapMonth(i); //闰哪个月
    this.isLeap = false;
    for (i = 1; i < 13 && offset > 0; i++) {
        if (leap > 0 && i == (leap + 1) && this.isLeap == false) {    //闰月
            --i;
            this.isLeap = true;
            temp = leapDays(this.year);
        }
        else {
            temp = monthDays(this.year, i);
        }
        if (this.isLeap == true && i == (leap + 1)) this.isLeap = false;    //解除闰月
        offset -= temp;
        if (this.isLeap == false) this.monCyl++;
    }
    if (offset == 0 && leap > 0 && i == leap + 1)
        if (this.isLeap) {
            this.isLeap = false;
        }
        else {
            this.isLeap = true;
            --i;
            --this.monCyl;
        }
    if (offset < 0) {
        offset += temp;
        --i;
        --this.monCyl;
    }
    this.month = i;
    this.day = offset + 1;
}

//返回公历y年m+1月的天数
function solarDays(y, m) {
    if (m == 1)
        return (((y % 4 == 0) && (y % 100 != 0) || (y % 400 == 0)) ? 29 : 28);
    else
        return (solarMonth[m]);
}

//记录公历和农历某天的日期
function calElement(sYear, sMonth, sDay, week, lYear, lMonth, lDay, isLeap) {
    this.isToday = false;
    //公历
    this.sYear = sYear;
    this.sMonth = sMonth;
    this.sDay = sDay;
    this.week = week;
    if (this.sYear === tY && this.sMonth === (tM + 1) && this.sDay === tD) {
        this.isToday = true;
    }
    //农历
    this.lYear = lYear;
    this.lMonth = lMonth;
    this.lDay = lDay;
    this.isLeap = isLeap;
    //节日记录
    this.lunarFestival = ''; //农历节日
    this.solarFestival = ''; //公历节日
    this.solarTerms = ''; //节气
}

//返回某年的第n个节气为几日(从0小寒起算)
function sTerm(y, n) {
    var offDate = new Date((31556925974.7 * (y - 1900) + sTermInfo[n] * 60000) + Date.UTC(1900, 0, 6, 2, 5));
    return (offDate.getUTCDate())
}

//保存y年m+1月的相关信息
var fat = mat = 9;
var eve = 0;

function calendar(y, m) {
    fat = mat = 0;
    var sDObj, lDObj, lY, lM, lD = 1, lL, lX = 0, tmp1, tmp2;
    var lDPOS = new Array(3);
    var n = 0;
    var firstLM = 0;
    sDObj = new Date(y, m, 1);    //当月第一天的日期
    this.length = solarDays(y, m);    //公历当月天数
    this.firstWeek = sDObj.getDay();    //公历当月1日星期几
    if ((m + 1) == 5) {
        fat = sDObj.getDay()
    }
    if ((m + 1) == 6) {
        mat = sDObj.getDay()
    }
    for (var i = 0; i < this.length; i++) {
        if (lD > lX) {
            sDObj = new Date(y, m, i + 1);    //当月第一天的日期
            lDObj = new Dianaday(sDObj.getTime());     //农历
            // console.log(sDObj,lDObj);
            lY = lDObj.year;           //农历年
            lM = lDObj.month;          //农历月
            lD = lDObj.day;            //农历日
            lL = lDObj.isLeap;         //农历是否闰月

            lX = lL ? leapDays(lY) : monthDays(lY, lM); //农历当月最后一天
            if (lM == 12) {
                eve = lX
            }
            if (n == 0) firstLM = lM;
            lDPOS[n++] = i - lD + 1;
        }
        // console.log(sDObj,lDObj);
        this[i] = new calElement(y, m + 1, i + 1, nStr1[(i + this.firstWeek) % 7], lY, lM, lD++, lL);
        // console.log(this[i]);
        if ((i + this.firstWeek) % 7 == 0) {
            this[i].color = 'red';  //周日颜色
        }
    }
    //节气
    tmp1 = sTerm(y, m * 2) - 1;
    tmp2 = sTerm(y, m * 2 + 1) - 1;
    this[tmp1].solarTerms = solarTerm[m * 2];
    this[tmp2].solarTerms = solarTerm[m * 2 + 1];
    if ((this.firstWeek + 12) % 7 == 5) {//黑色星期五
        this[12].solarFestival += '黑色星期五';
    }
    if (y == tY && m == tM) {
        this[tD - 1].isToday = true;    //今日
    }
}

function calendarWeek(weekStartDay) {
    var fat = mat = 0;
    var sDObj, lDObj, lY, lM, lD = 1, lL, lX = 0, tmp1, tmp2, tmp3, tmp4;
    var lDPOS = new Array(3);
    var n = 0;
    var firstLM = 0;
    //02-08
    sDObj = new Date(weekStartDay);    //这周第一天的日期
    var y = sDObj.getFullYear();
    var m = sDObj.getMonth() + 1;
    // var d = sDObj.getDate();
    var endDateTime = new Date(weekStartDay).getTime() + 6 * 24 * 60 * 60 * 1000;
    var endDate = new Date(endDateTime);
    // var endDateFormat = formatTime(endDateTime);
    // var endY = endDate.getFullYear();
    var endM = endDate.getMonth() + 1;

    this.length = 7;
    this.firstWeek = sDObj.getDay();
    if (m == 5) {
        fat = sDObj.getDay()
    }
    if (m == 6) {
        mat = sDObj.getDay()
    }

    for (var i = 0; i < this.length; i++) {
        var dateTime = new Date(weekStartDay).getTime() + i * 24 * 60 * 60 * 1000;
        var currentDate = new Date(dateTime);
        // console.log('lX-->' + lX);
        // console.log('lD-->' + lD);

        if (lD > lX) {
            sDObj = new Date(weekStartDay);
            // lDObj = new Dianaday(sDObj);     //农历
            lDObj = new Dianaday(dateTime);     //农历
            // console.log(sDObj,lDObj);
            lY = lDObj.year;           //农历年
            lM = lDObj.month;          //农历月
            lD = lDObj.day;            //农历日
            lL = lDObj.isLeap;         //农历是否闰月

            lX = lL ? leapDays(lY) : monthDays(lY, lM); //农历当月最后一天
            if (lM == 12) {
                eve = lX
            }
            if (n == 0) firstLM = lM;
            lDPOS[n++] = i - lD + 1;
        }
        // console.log(lDObj)

        var indexName = getIndexName(currentDate.getMonth(), currentDate.getDate());
        this[indexName] = new calElement(currentDate.getFullYear(), currentDate.getMonth() + 1, currentDate.getDate(), nStr1[(i + this.firstWeek) % 7], lY, lM,  lD++, lL);

    }


    //节气
    tmp1 = sTerm(y, m * 2) - 1;
    tmp2 = sTerm(y, m * 2 + 1) - 1;
    tmp3 = sTerm(y, endM * 2) - 1;
    tmp4 = sTerm(y, endM * 2 + 1) - 1;
    this[getIndexName(m, tmp1)] && (this[getIndexName(m, tmp1)].solarTerms = solarTerm[m * 2]);
    this[getIndexName(m, tmp2)] && (this[getIndexName(m, tmp2)].solarTerms = solarTerm[m * 2 + 1]);
    this[getIndexName(endM, tmp3)] && (this[getIndexName(endM, tmp3)].solarTerms = solarTerm[endM * 2]);
    this[getIndexName(endM, tmp4)] && (this[getIndexName(endM, tmp4)].solarTerms = solarTerm[endM * 2 + 1]);
    // this[tmp2] && (this[tmp2].solarTerms = solarTerm[m * 2 + 1]);
    // if ((this.firstWeek + 12) % 7 == 5) {//黑色星期五
    //     this[12].solarFestival += '黑色星期五';
    // }
    var toDayName = (tM+1).toString().padStart(2, '00') + '-' + tD.toString().padStart(2, '00')
    this[toDayName] && (this[toDayName]['isToday'] = true);
}

function getIndexName(m, d) {
    return (m + 1).toString().padStart(2, '00') + '-' + d.toString().padStart(2, '00')
}

function formatTime(time) {
    let oDate = new Date(time * 1);
    let oYear = oDate.getFullYear();
    let oMonth = oDate.getMonth() + 1;
    let oDay = oDate.getDate();
    return oYear + '-' + oMonth.toString().padStart(2, '00') + '-' + oDay.toString().padStart(2, '00');
}