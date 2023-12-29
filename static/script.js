import  {renderEmployee, addDepartment} from "/static/renderforms.js"

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
            disable();
        });
    })
    
    const deleteButton = document.querySelector('#delete');
    deleteButton.addEventListener(`click`, ()=>{
        //TODO: Add a DELETE request (for later)
    })

    const newItem = document.querySelector('.new-item');
    const newItemDrop = newItem.querySelector('.dropdown-box');

    const employeeBox = newItemDrop.children[0];
    const departmentBox = newItemDrop.children[1];

    employeeBox.addEventListener(`click`, ()=>{
        renderEmployee(employee, 'POST');
    });

    const updateButton = document.querySelector('#update');
    updateButton.addEventListener('click', ()=>{
        const select = document.querySelector('.selected');

        const emp = convertFieldToEmployee(select);
        renderEmployee(emp, 'PUT');
        

    });
    departmentBox.addEventListener(`click`, addDepartment);

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

function disable(){
    const len = [...document.querySelectorAll('#select')].filter(i => i.checked).length;
    const maxLen = [...document.querySelectorAll('#select')].length;
    const selectAllButton = document.querySelector('#select-all');

    if(len != maxLen) selectAllButton.checked = false;
    
}

const employee = (id, first, last, email, department, departmentId) => {
    return {id, first, last, email, department, departmentId}
}

//Get employee info from an HTML element
const convertFieldToEmployee = (select) => {

        const id = Number(select.querySelector("#employee-id").textContent);
        const lastName = select.querySelector('#last-name').textContent;
        const firstName = document.querySelector('#first-name').textContent;
        const email = document.querySelector('#email-data').textContent;
        const dept = document.querySelector('[data-dept-id]').textContent;
        const deptId = Number(document.querySelector('[data-dept-id]').dataset.deptId);

        const emp = employee(id, firstName, lastName, email, dept, deptId);
        return emp;
}

main();