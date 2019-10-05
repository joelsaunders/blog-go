import React from 'react';
import {ErrorMessage, Field} from "formik";
import {Link} from "react-router-dom";

const PageOne = (props) => {
    const titleErr = props.formProps.errors.hasOwnProperty('title');
    const descriptionErr = props.formProps.errors.hasOwnProperty('description');
    const pictureErr = props.formProps.errors.hasOwnProperty('picture');

    return <>
        <div className="my-8">
            <label className="block text-gray-700 text-base font-bold mb-2" htmlFor="title">
                Title
            </label>
            <Field type="text"
                   name="title"
                   placeholder="Title"
                   className="shadow appearance-none border rounded w-full sm:w-1/2 py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            />
        </div>
        <div className={`ui error message ${titleErr ? 'visible': null}`}>
            <ErrorMessage className="ui error message visible" name="title" />
        </div>

        <div className="my-4">
            <label className="block text-gray-700 text-base font-bold mb-2" htmlFor="description">
                Description
            </label>
            <Field component="textarea"
                   rows="10"
                   name="description"
                   placeholder="Description"
                   className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            />
        </div>
        <div className={`ui error message ${descriptionErr ? 'visible': null}`}>
            <ErrorMessage className="ui error message visible" name="description" />
        </div>

        <div className="my-4">
            <label className="block text-gray-700 text-base font-bold mb-2" htmlFor="description">
                Picture
            </label>
            <Field type="text"
                   name="picture"
                   placeholder="Picture"
                   className="shadow appearance-none border rounded w-full sm:w-1/2 py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            />
        </div>
        <div className={`ui error message ${pictureErr ? 'visible': null}`}>
            <ErrorMessage className="ui error message visible" name="picture" />
        </div>

        <Link to={`${props.rootURL}/page-two`}>
            <div className="inline-block text-sm px-4 py-2 leading-none border rounded text-teal-500 border-white hover:border-transparent hover:text-white hover:bg-teal-500 mt-4 mb-4">
                Next Page
            </div>
        </Link>
    </>
};

export default PageOne;
