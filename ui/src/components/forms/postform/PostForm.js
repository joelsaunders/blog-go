import React from 'react';
import {Form, withFormik} from "formik";
import * as Yup from 'yup';
import {Redirect, Route, Switch} from "react-router-dom";
import PageOne from "./PageOne";
import PageTwo from "./PageTwo";
import customHistory from "../../../customHistory";

const PostForm = (props) => {
    return <div>
        <Form>
            <Switch>
            <Redirect from={props.match.path} exact to={`${props.match.path}/page-one`}/>
            <Route
                path={`${props.match.path}/:page-one`}
                render={routeProps => (
                        <PageOne
                            {...routeProps}
                            formProps={props}
                            rootURL={props.match.url}
                        />
                    )}
            />
            <Route
                path={`${props.match.path}/:page-two`}
                render={routeProps => (
                        <PageTwo
                            {...routeProps}
                            formProps={props}
                            rootURL={props.match.url}
                        />
                    )}
            />
            </Switch>
        </Form>
    </div>
};

const FormikPostForm = withFormik({
    mapPropsToValues: (props) => {
        return {
            title: props.post? props.post.title || '': '',
            slug: props.post? props.post.slug || '': '',
            description: props.post? props.post.description || '': '',
            body: props.post? props.post.body || '': '',
            picture: props.post? props.post.picture || '': ''
        }
    },
    validationSchema: Yup.object().shape({
        title: Yup.string().max(100, 'Must be less than 100 characters').required('Required'),
        slug: Yup.string().max(50, 'Must be less than 50 characters').required('Required')
    }),
    enableReinitialize: true,
    handleSubmit: (values, {setSubmitting, props}) => {
        props.submitAction(values);
        setSubmitting(false);
        customHistory.push("/")
    }
});

export default FormikPostForm(PostForm);