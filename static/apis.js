/*
    Methods here are responsible for API methods
*/

import { renderNewEmployee } from "./script.js";

//Master method for methods
export async function submitEmployee(employee, Method){
    console.log(employee);
    switch(Method){
        case 'PUT':
            const updatedEmp = await putEmployee(employee);
            break;
        case 'POST':
            const newEmp = await postEmployee(employee)
            //Hand it over to the table for new rendering
            console.log(`Employee: ${newEmp}`);
            renderNewEmployee(newEmp)
            break;
    }
}

export async function deleteEmployee(){
    
}

async function postEmployee(employee){
    try{
        //console.log(employee);
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

}