/*
    Methods here are responsible for API methods
*/

import { renderNewEmployee, renderExistingEmployee, removeDeletedEmployees } from "./script.js";

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
            renderNewEmployee(newEmp)
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
        removeDeletedEmployees();
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
        const emp = await response.json();

        //post employee for the frontend
        return emp;
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
        const updatedEmp = await response.json();
        return updatedEmp
    }catch(e){
        console.error(e);
    }
}