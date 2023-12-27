

function main(){
    const selectAllButton = document.querySelector('#select-all');
    selectAllButton.addEventListener('click', (e)=>{
       checkAll();
    })
}

function checkAll(){
    [...document.querySelectorAll('#select')].forEach(ele => {
        ele.checked = !ele.checked;
        console.log(ele.parentNode.parentNode);
    })
}
main();