import React, {useEffect} from 'react';
import {Form, withFormik} from "formik";
import * as Yup from 'yup';
import {Redirect, Route, Switch} from "react-router-dom";
import PageOne from "./PageOne";
import PageTwo from "./PageTwo";
import customHistory from "../../../customHistory";
import {connect} from "react-redux";
import {fetchTags} from "../../../actions";

const PostForm = (props) => {
    const fetchTags = props.fetchTags;
    useEffect(
        () => {
            (async () =>  await fetchTags())()
        },
        [fetchTags]
    );

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
            picture: props.post? props.post.picture || '': '',
            tags: props.post? props.post.tags || []: [],
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

const mapStateToProps = (state) => {
    return {tags: state.tags}
};

export default connect(mapStateToProps, {fetchTags})(FormikPostForm(PostForm));