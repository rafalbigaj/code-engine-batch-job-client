/*******************************************************************************
 * Licensed Materials - Property of IBM
 * IBM Cloud Code Engine, 5900-AB0
 * © Copyright IBM Corp. 2020
 * US Government Users Restricted Rights - Use, duplication or
 * disclosure restricted by GSA ADP Schedule Contract with IBM Corp.
 ******************************************************************************/

package pointers

func Int64Ptr(i int64) *int64 { return &i }

func StringPtr(s string) *string { return &s }

func BoolPtr(b bool) *bool { return &b }
