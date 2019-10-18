import React from 'react';
import theBookOfJoel from "../../apis/theBookOfJoel";
import {Field, Form, Formik} from "formik";
import {connect} from "react-redux";
import {signOut} from "../../actions";

const ChangePasswordForm = (props) => {
    const handleSubmit = async (values, actions) => {
        let response;
        try {
            response = await theBookOfJoel.patch(
                `api/v1/user/change-password`,
                values,
                {headers: {Authorization: `Bearer ${props.currentUserToken}`}}
            );
        } catch (e) {
            if (e.response.status === 400) {
                actions.setErrors(e.response.data);
            }
            actions.setSubmitting(false);
            actions.setStatus({ msg: JSON.stringify(e.response.data)});
            return
        }
        actions.setSubmitting(false);
        props.signOut();
        props.onDismiss();
        //TODO: potentially open login modal here
        return response
    };

    return <div>
        <Formik
            initialValues={{password:""}}
            onSubmit={(values, actions) => handleSubmit(values, actions)}
        >
            { props => {
                const { status, isSubmitting} = props;

                return (
                    <Form>
                        <div className="mb-4">
                            <label
                                className="block text-gray-700 text-sm font-bold mb-2"
                                htmlFor="password"
                            >
                                Password
                            </label>
                            <Field
                                className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                                type="password"
                                name="password"
                            />
                            {status && status.msg && <div>{status.msg}</div>}
                        </div>
                        <button
                            className="inline-block text-sm px-4 py-2 leading-none border rounded text-teal-500 border-white hover:border-transparent hover:text-white hover:bg-teal-500 mt-4"
                            type="submit"
                            disabled={isSubmitting}
                        >
                            Submit
                        </button>
                    </Form>
                );
            }}
        </Formik>
    </div>
};

const mapStateToProps = (state) => {
    return {currentUserToken: state.auth.token}
};

export default connect(mapStateToProps, {signOut})(ChangePasswordForm);
