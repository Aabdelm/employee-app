import  {addEmployee} from "/static/renderforms.js"

//Responsible for dealing with the DOM
function main(){
    const selectAllButton = document.querySelector('#select-all');
    selectAllButton.addEventListener('click', (e)=>{
       checkAll(e.target.checked);
    })

    const select = document.querySelectorAll('#select');
    select.forEach(ele =>{
        ele.addEventListener(`click`, (e)=>{
            //Somewhat lazy approach here
            controlVisibility(e.target);
            //in case the checkall button was previously enabled
            disableTest();
        });
    })
    
    const deleteButton = document.querySelector('#delete');
    deleteButton.addEventListener(`click`, ()=>{
        //TODO: Add a DELETE request (for later)
    })

    const newItem = document.querySelector('.new-item');
    const newItemDrop = newItem.querySelector('.dropdown-box');
    const employeeBox = newItemDrop.children[0];

    employeeBox.addEventListener(`click`,addEmployee);

}




//Checks the button and enables the delete button accordingly.
function checkAll(checkedButton){
    [...document.querySelectorAll('#select')].forEach(ele => {
        ele.checked = checkedButton;
        controlVisibility(ele)
    })
}

function controlVisibility(ele){
    //Need to indicate it's selected
    const grandParent = ele.parentElement.parentElement;
    ele.checked ? grandParent.classList.add('selected') : grandParent.classList.remove('selected');
    //get the disable buton
    const deleteButton = document.querySelector('#delete');
    const updateButton = document.querySelector('#update');

    const len = [...document.querySelectorAll('#select')].filter(i => i.checked).length;

    deleteButton.style.visibility = len ? 'visible' : 'hidden'; 

    //We'll need to update only one department
    updateButton.style.visibility = len == 1 ? 'visible' : 'hidden'; 
    
}

function disableTest(){
    const len = [...document.querySelectorAll('#select')].filter(i => i.checked).length;
    const maxLen = [...document.querySelectorAll('#select')].length;
    const selectAllButton = document.querySelector('#select-all');

    if(len != maxLen) selectAllButton.checked = false;
    
}
main();