var socket = io.connect();

/*
let socket = new WebSocket("ws://localhost:1332/ws")
console.log("zkouším websocket")
socket.onopen = () => {
    console.log("Připojeno")
    socket.send("Zdar, tady klient")
}

socket.send*/




$(document).ready(function(){

    
    $("#socket").click(function(){  
        console.log("L")
        socket.emit("getImage","kokokoko","bla","kokokoko2","bla2")
    })

    $('#xhrContent').on('click','#imgList tr td a',function(){
        var imgName = $(this).attr("rel");
        socket.emit("get-image",imgName);
        return false;
	})

    socket.on("image", (img) => {
        data = JSON.parse(img)
        $("#imgContainer").html('<img src="data:image/'+data.imgType+';base64,'+data.imgBase64Content+'" />')
    });










    $('#xhrContentzal').on('click','#imgList tr td a',function(){
        var imgName = $(this).attr("rel");
        $.getJSON("/get-img/"+imgName,function(data){
            console.log(data)
           $("#imgContainer").html('<img src="data:image/'+data.imgType+';base64,'+data.imgBase64Content+'" />')
        });

        return false;
	})





    $("#imgUpload").submit(function(){
        var form = document.querySelector("#imgUpload")
        var formData = new FormData(form)
        for(let[name,value] of formData){
            alert(`${name} = ${value}`); 
        }


        $.ajax({
            url: "/img-upload",
            type: "POST",
            data: formData,
            success : function(data){
                var str = "";
                var json = $.parseJSON(data)
                if(json.imgName.length > 0){
                    str += '<div style="margin-top:3px" id="imgContainer"></div>';
                    str += '<table id="imgList">';
                    for(var i=0;i < json.imgName.length;i++){
                        str += '<tr><td><a href="javascript: void(0);" rel="'+json.imgName[i]+'">'+json.imgName[i]+'</a><div class="imgContainer"></div></td></tr>';
                    }
                    str += '</table>';
                }
                $("#xhrContent").html(str); 
              
/*

                $("#bla").html('<img src="data:image/'+json.imgType+';base64,'+json.imgBase64Content+'" />') */
            },
            enctype: 'multipart/form-data',
            processData: false,
            contentType: false,
            cache: false
        });
	    return false;
	})




    


});