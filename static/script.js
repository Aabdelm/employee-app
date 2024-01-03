import  {renderEmployee, addDepartment} from "/static/renderforms.js"
import { deleteEmployees, query} from "./apis.js";

//Responsible for dealing with the DOM
function main(){
    const selectAllButton = document.querySelector('#select-all');
    selectAllButton.addEventListener('click', (e)=>{
       checkAll(e.target.checked);
    })

    const select = document.querySelectorAll('#select');
    select.forEach(ele =>{
        ele.addEventListener(`click`, renderVisibility);
    })
    
    const deleteButton = document.querySelector('#delete');
    deleteButton.addEventListener(`click`, ()=>{
        const emps = [];
        const selectedAll = document.querySelectorAll('.selected');
        selectedAll.forEach(ele => {
            const emp = convertFieldToEmployee(ele);
            emps.push(emp);
        })
        deleteEmployees(emps);
    })

    const newItem = document.querySelector('.new-item');
    const newItemDrop = newItem.querySelector('.dropdown-box');

    const employeeBox = newItemDrop.children[0];
    const departmentBox = newItemDrop.children[1];

    employeeBox.addEventListener(`click`, ()=>{
        renderEmployee(Employee, 'POST');
    });

    const updateButton = document.querySelector('#update');
    updateButton.addEventListener('click', ()=>{
        const select = document.querySelector('.selected');

        const emp = convertFieldToEmployee(select);
        renderEmployee(emp, 'PUT');
        

    });
    departmentBox.addEventListener(`click`, addDepartment);

    const search = document.querySelector('[type=search]');
    search.addEventListener(`keydown`, (e)=>{
        if(e.key == 'Enter' || e.key == 'Return'){
            e.preventDefault();
            const option = [...document.querySelectorAll('option')]
            .filter(element => element.selected)[0];

            const val = option.value;
            const identifer = search.value;

            query(identifer, val)
        }
    })

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

export function Employee(id, firstName, lastName, email, department, departmentId){
    return {id, firstName, lastName, email, department, departmentId}
}


export async function renderNewEmployee(emp, method){
    //debugger;
    const tbody = document.querySelector('tbody');
    const tr = document.createElement('tr');
    const empBox = document.querySelectorAll('.figure')[0];
    const empCounter = counter(empBox.querySelector('.number'));
    const info = empBox.querySelector('.info');

    const buttonParent = document.createElement('td');
    const button = document.createElement('input');
    button.id = 'select';
    button.type = 'checkbox';
    button.name = 'checkbox';

    //Add button functionality
    button.addEventListener('click', renderVisibility);

    buttonParent.appendChild(button);
    tr.appendChild(buttonParent);

    const empId = document.createElement('td');
    empId.id = 'employee-id';
    empId.textContent = emp.id;
    tr.appendChild(empId);

    const lastName = document.createElement('td');
    lastName.id = 'last-name';
    lastName.textContent = emp.lastName;
    tr.appendChild(lastName);

    const firstName = document.createElement('td');
    firstName.id = 'first-name';
    firstName.textContent = emp.firstName;
    tr.appendChild(firstName);

    const email = document.createElement('td');
    email.id = 'email-data';
    email.textContent = emp.email;
    tr.appendChild(email);

    const dept = document.createElement('td');
    dept.dataset['deptId'] = emp.departmentId;
    dept.textContent = emp.department;
    tr.appendChild(dept);

    tbody.appendChild(tr);

    if(method == "POST") empCounter.increaseCount()
    renderPlurality(info,'Employee',empCounter);

    
    
}

export function clearAndReRender(emps){
    const rows = [...document.querySelectorAll('tr')]
    .filter((_, index) => index > 0)
    .forEach(element => element.remove())

    emps.forEach(emp => {
        renderNewEmployee(emp);
    })
     
}

export function toastify(err){
    const toast = document.querySelector('.toast');
    const error = toast.querySelector('.error');
    error.textContent = err;

    toast.style.transition = "0.3s ease";
    setTimeout(()=>{
        toast.style.opacity = '100%'
    }, 500);
    setTimeout(()=>{
        toast.style.opacity = '0';
    },3500);
        
}

export async function renderExistingEmployee(emp){
    const select = document.querySelector('.selected');

    const lastName = select.querySelector('#last-name');
    lastName.textContent = emp.lastName;

    const firstName = select.querySelector('#first-name');
    firstName.textContent = emp.firstName;

    const email = select.querySelector('#email-data');
    email.textContent = emp.email;

    const dept = select.querySelector('[data-dept-id]');
    dept.textContent = emp.department;

    dept.dataset.deptId = emp.departmentId;
}

const counter = (ele) => {
    let count = Number(ele.textContent);
    const increaseCount = ()=>{
        count++;
        ele.textContent = count;
    }
    const decreaseCount = () => {
        count--;
        ele.textContent = count;
    }
    const getCount = () => {
        return count;
    }
    return {getCount, increaseCount, decreaseCount}
}

export function removeDeletedEmployees(){
    const selected = document.querySelectorAll('.selected');

    const empBox = document.querySelectorAll('.figure')[0]
    const figures = empBox.querySelector('.info');
    const empCounter = counter(empBox.querySelector('.number'));
    selected.forEach(ele => {
        ele.remove();
        empCounter.decreaseCount();
        renderPlurality(figures, 'Employee', empCounter);
    })

    document.querySelector('#select-all').checked = false;
}

export function renderDepartmentAddition(){
    const deptBox = document.querySelectorAll('.figure')[1];
    const deptCounter = counter(deptBox.querySelector('.number'));
    const info = deptBox.querySelector('.info');

    deptCounter.increaseCount();
    renderPlurality(info, 'Department', deptCounter);
}


function renderPlurality(element, word,counter){
    element.textContent = counter.getCount() == 1 ? word : (word + 's');
}

//Get employee info from an HTML element
const convertFieldToEmployee = (select) => {
        const id = Number(select.querySelector("#employee-id").textContent);
        const lastName = select.querySelector('#last-name').textContent;
        const firstName = select.querySelector('#first-name').textContent;
        const email = select.querySelector('#email-data').textContent;
        const dept = select.querySelector('[data-dept-id]').textContent;
        const deptId = Number(select.querySelector('[data-dept-id]').dataset.deptId);

        const emp = Employee(id, firstName, lastName, email, dept, deptId);
        return emp;
}

function renderVisibility(e){
            //Somewhat lazy approach here
            controlVisibility(e.target);
            //in case the checkall button was previously enabled
            disable();
}

main();