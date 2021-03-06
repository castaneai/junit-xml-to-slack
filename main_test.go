package main

import (
	"testing"
	"strings"
	"encoding/json"
	"bytes"
)

func TestConvert(t *testing.T) {
	input := `<?xml version="1.0" ?>
<testsuite name="info.classmethod.testing.HogeTest" tests="4" errors="1" failures="1" time="0.003">
  <testcase classname="info.classmethod.testing.HogeTest" name="testCase01" time="0.001" />
  <testcase classname="info.classmethod.testing.HogeTest" name="testCase02" time="0.001">
    <failure type="java.lang.AssertionError" message="test failed" >java.lang.AssertionError: test failed
  at org.junit.Assert.fail(Assert.java:88)
  at jp.classmethod.testing.examples.core.runner.DynamicTestsRunnerExample$1.invokeTest(DynamicTestsRunnerExample.java:33)
  at jp.classmethod.testing.core.runner.DynamicTestsRunner.invokeTest(DynamicTestsRunner.java:86)
  at jp.classmethod.testing.core.runner.DynamicTestsRunner.run(DynamicTestsRunner.java:68)
  at org.eclipse.jdt.internal.junit4.runner.JUnit4TestReference.run(JUnit4TestReference.java:50)
  at org.eclipse.jdt.internal.junit.runner.TestExecution.run(TestExecution.java:38)
  at org.eclipse.jdt.internal.junit.runner.RemoteTestRunner.runTests(RemoteTestRunner.java:467)
  at org.eclipse.jdt.internal.junit.runner.RemoteTestRunner.runTests(RemoteTestRunner.java:683)
  at org.eclipse.jdt.internal.junit.runner.RemoteTestRunner.run(RemoteTestRunner.java:390)
  at org.eclipse.jdt.internal.junit.runner.RemoteTestRunner.main(RemoteTestRunner.java:197)</failure>
  </testcase>
  <testcase classname="info.classmethod.testing.HogeTest" name="testCase03" time="0.001">
    <error type="java.lang.NullPointerException" message="NPE" >java.lang.NullPointerException: NPE
  at jp.classmethod.testing.examples.core.runner.DynamicTestsRunnerExample$1.invokeTest(DynamicTestsRunnerExample.java:36)
  at jp.classmethod.testing.core.runner.DynamicTestsRunner.invokeTest(DynamicTestsRunner.java:86)
  at jp.classmethod.testing.core.runner.DynamicTestsRunner.run(DynamicTestsRunner.java:68)
  at org.eclipse.jdt.internal.junit4.runner.JUnit4TestReference.run(JUnit4TestReference.java:50)
  at org.eclipse.jdt.internal.junit.runner.TestExecution.run(TestExecution.java:38)
  at org.eclipse.jdt.internal.junit.runner.RemoteTestRunner.runTests(RemoteTestRunner.java:467)
  at org.eclipse.jdt.internal.junit.runner.RemoteTestRunner.runTests(RemoteTestRunner.java:683)
  at org.eclipse.jdt.internal.junit.runner.RemoteTestRunner.run(RemoteTestRunner.java:390)
  at org.eclipse.jdt.internal.junit.runner.RemoteTestRunner.main(RemoteTestRunner.java:197)</error>
  </testcase>
  <testcase classname="info.classmethod.testing.HogeTest" name="testCase04" time="0.000">
    <skipped/>
  </testcase>
  <system-out><![CDATA[stdout!]]></system-out>
  <system-err><![CDATA[stderr!]]></system-err>
</testsuite>`

	jr := strings.NewReader(input)
	suite, err := parseJUnitXML(jr)
	if err != nil {
		t.Fatal(err)
	}

	payload, err := toSlackPayload(suite)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(payload); err != nil {
		t.Fatal(err)
	}

	t.Log(buf.String())
}

