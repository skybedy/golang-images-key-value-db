$(document).ready(function(){
    var socket = io.connect();

    $("#file").change(function(){
        const reader = new FileReader()
        reader.fileName = this.files[0].name
        fileNameArr = reader.fileName.split(".")
        fileType = fileNameArr[fileNameArr.length-1]
        if(fileType !== "jpg" && fileType !== "gif" && fileType !== "png" && fileType !== "JPG" && fileType !== "GIF" && fileType !== "PNG"){
            alert("Je možné nahrát jen obrázky ve formátu jpg, gif nebo png")
            return false;
        }

        reader.onload = function(readerEvt){
            const base64 = this.result.replace(/.*base64,/, '');
            socket.emit('images', base64,readerEvt.target.fileName); 
        }
        reader.readAsDataURL(this.files[0]);
        return false;
    })

    $("#file").click(function(){
        $("#imgContainer").html(""); 
    })

    socket.on("images", (img) => {
        json = JSON.parse(img);
        var str = "";
        str += '<table>';
        for(var i=0;i < json.imgName.length;i++){
            str += '<tr><td><a href="javascript: void(0);" rel="'+json.imgName[i]+'">'+json.imgName[i]+'</a></td></tr>';
        }
        str += '</table>';
        $("#imgList").html(str); 
    });


    $('#xhrContent').on('click','#imgList tr td a',function(){
        var imgName = $(this).attr("rel");
        socket.emit("get-image",imgName);
        return false;
	})

    socket.on("image", (img) => {
        data = JSON.parse(img)
        $("#imgContainer").html('<img src="data:image/'+data.imgType+';base64,'+data.imgBase64Content+'" />')
    });

    $("#truncateDbpg").click(function(){
        $("#imgContainer").html(""); 
        socket.emit("truncatedbpg");
        $("#imgList").html(""); 
    })

});
