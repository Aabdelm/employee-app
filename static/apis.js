/*
    Methods here are responsible for API methods
*/

import { renderNewEmployee, renderExistingEmployee, 
    removeDeletedEmployees, renderDepartmentAddition, clearAndReRender, toastify } from "./script.js";

//Master method for methods
export async function submitEmployee(employee, Method){
    switch(Method){
        case 'PUT':
            const updatedEmp = await putEmployee(employee);
            renderExistingEmployee(updatedEmp);
            break;
        case 'POST':
            const newEmp = await postEmployee(employee)
            //Hand it over to the table for new rendering
            renderNewEmployee(newEmp, "POST")
            break;
    }
}

export async function deleteEmployees(employees){
    try{
        const response = await fetch('http://localhost:80/employees/',{
            method: 'DELETE',
            headers:{
                "Content-Type": "application/JSON",
            },
            body: JSON.stringify(employees),
        });
        if(!response.ok){
            awaitAndToastify(response);
        }
        else removeDeletedEmployees();
    }catch(e){
        console.error(e);
    }
}

async function postEmployee(employee){
    try{
        const response = await fetch('http://localhost:80/employees/',{
            method: 'POST',
            headers:{
                "Content-Type": "application/JSON",
            },
            body: JSON.stringify(employee)
        });
        if(!response.ok){
            awaitAndToastify(response);
        }else{
            const emp = await response.json();
            //post employee for the frontend
            return emp;
        }

        
    }catch(e){
        console.error(e);
    }
}

async function putEmployee(employee){
    try{
        const response = await fetch(`http://localhost:80/employees/${employee.id}`,{
        method: 'PUT',
        headers:{
            "Content-type": "application/JSON",
        },
        body: JSON.stringify(employee),
    });
        if(!resp.ok){
            awaitAndToastify(response);
        }else{
            const updatedEmp = await response.json();
            return updatedEmp
        }
    }catch(e){
        console.error(e);
    }
}

export async function submitDepartment(dept){
    try{
        const res = await fetch(`http://localhost:80/departments/`,{
            method: 'POST',
            headers:{
                "Content-type": "application/JSON",
            },
            body: JSON.stringify(dept),
        })
        if(!res.ok){
            awaitAndToastify(res);
        }
        else renderDepartmentAddition();
    }catch(e){
        console.error(e)
    }
}

export async function query(identifier, value){
    
    const query = value.split("").filter(char => char != '-').join("");
    console.log(value);
    if(value == 'employee-id' && !isDigit(identifier)) return;

    try{
        const resp = value !== 'employee-id' ? 
        await fetch(`http://localhost:80/employees?${query}=${identifier}`) : 
        await fetch(`http://localhost:80/employees/${identifier}`)
        if (!resp.ok){
            awaitAndToastify(resp);
        }else{
            const emps = value !== 'employee-id' ? [...await resp.json()] : [await resp.json()];
            clearAndReRender(emps);
        }
        
    }catch(e){
        console.error(e)
    }
}

const isDigit = (digit) =>{
    const regEx = new RegExp("^[0-9]+$");
    return regEx.test(digit);
}

export async function awaitAndToastify(res){
    const err = await res.text();
    toastify(err);
}