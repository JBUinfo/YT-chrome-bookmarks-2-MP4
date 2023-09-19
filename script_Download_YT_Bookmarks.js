const fs = require('fs');
const { exec } = require( 'child_process' );
let jsonData = require('./download.json');
let arrayDownloaded = fs.existsSync('downloaded.txt') ? fs.readFileSync('downloaded.txt').toString().split("\n") : [];
let downloaded = fs.createWriteStream("downloaded.txt", {flags:'a'});
let errors = fs.createWriteStream("errors.txt", {flags:'a'});

function sleep(ms){ //Function Sleep
  return new Promise(resolve=>{
    setTimeout(resolve,ms)
  })
}

let cont = 0;
let contVideosTotales = 0;
let contVideosDownloaded = 0;
async function recursiveFolders(node, path){
  contVideosTotales += node.length;
	for(let i= 0; i < node.length;i++){
		while(cont == 3){await sleep(1000);}
    		if(node[i].type == "folder") recursiveFolders(node[i].items, path+"/"+node[i].title)
    		if(node[i].type == "link" && node[i].href.includes("youtube.com")) createFile(node[i],path)

	}
}

function createFile(item,path){
  cont++;
  console.log(path, item.title);
  if (!fs.existsSync(path)) fs.mkdirSync(path, { recursive: true });
  let bo = false;
  for(i in arrayDownloaded) {
      if(i == item.href){
        bo = true;
        break;
      }
  }
  if (!bo) {
    exec(`"./yt-dlp.exe" -o "${path}/%(title)s.mp4" --restrict-filenames --no-check-certificate "${item.href}"`, (error, stdout, stderr) => {
        if (error){
          console.log("ERROR EN " + item.title + "\n");
	  console.log("ERROR:" + error.message + "\n");
          errors.write(error.message + "\n\n\n");
        } else if (stderr){
	  console.log("ERROR EN " + item.title + "\n");
	  console.log("ERROR:" + stderr + "\n");
          errors.write(stderr + "\n\n\n");
	} else {
          contVideosDownloaded++;
          console.log(contVideosDownloaded +"/"+ contVideosTotales);
          downloaded.write(item.href + "\n");
          arrayDownloaded.push(item.href);
        }
          cont--;
    });
  }
}
(async ()=>{
  await recursiveFolders(jsonData.folders[0].items[0].items, "./"+jsonData.folders[0].items[0].title);
  while (0 != cont){await sleep(1000);}
  downloaded.end();
  errors.end();
})()
