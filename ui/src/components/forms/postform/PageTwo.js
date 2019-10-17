import React from 'react';
import {ErrorMessage, Field} from "formik";


const PageTwo = (props) => {
    const err = props.formProps.errors.hasOwnProperty('description');

    return <>
        <div className="field">
            <label className="block text-gray-700 text-base font-bold mb-2" htmlFor="body">
                Description
            </label>
            <Field component="textarea"
                   rows="20"
                   name="body"
                   placeholder="Body"
                   className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            />
        </div>
        <div className={`ui error message ${err ? 'visible': null}`}>
            <ErrorMessage className="ui error message visible" name="body" />
        </div>
        <button
            className="inline-block text-sm px-4 py-2 leading-none border rounded text-teal-500 border-white hover:border-transparent hover:text-white hover:bg-teal-500 mt-4 mb-4"
            type="submit"
            disabled={props.formProps.isSubmitting}
        >
            Submit
        </button>
    </>
};

export default PageTwo;
