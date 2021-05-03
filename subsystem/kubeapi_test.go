package subsystem

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/jinzhu/gorm"
	bmhv1alpha1 "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/openshift/assisted-service/client"
	"github.com/openshift/assisted-service/internal/common"
	"github.com/openshift/assisted-service/internal/controller/api/v1beta1"
	"github.com/openshift/assisted-service/internal/controller/controllers"
	"github.com/openshift/assisted-service/internal/gencrypto"
	"github.com/openshift/assisted-service/models"
	"github.com/openshift/assisted-service/pkg/auth"
	conditionsv1 "github.com/openshift/custom-resource-status/conditions/v1"
	hivev1 "github.com/openshift/hive/apis/hive/v1"
	agentv1 "github.com/openshift/hive/apis/hive/v1/agent"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
)

var BASIC_KUBECONFIG = `apiVersion: v1
clusters:
- cluster:
    certificate-authority: /home/pawan/.minikube/ca.crt
    extensions:
    - extension:
        last-update: Wed, 28 Apr 2021 23:03:12 EDT
        provider: minikube.sigs.k8s.io
        version: v1.19.0
      name: cluster_info
    server: https://192.168.39.88:8443
  name: minikube
contexts:
- context:
    cluster: minikube
    extensions:
    - extension:
        last-update: Wed, 28 Apr 2021 23:03:12 EDT
        provider: minikube.sigs.k8s.io
        version: v1.19.0
      name: context_info
    namespace: default
    user: minikube
  name: minikube
current-context: minikube
kind: Config
preferences: {}
users:
- name: minikube
  user:
    client-certificate: /home/pawan/.minikube/profiles/minikube/client.crt
    client-key: /home/pawan/.minikube/profiles/minikube/client.key`
var BASIC_KUBECONFIG_OLD = `apiVersion: v1
clusters: 
  - 
    cluster: 
      certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURvekNDQW91Z0F3SUJBZ0lVTENPcVdURmFFQThnTkVtVityYjdoMXYwcjNFd0RRWUpLb1pJaHZjTkFRRUxCUUF3WVRFTE1Ba0dBMVVFQmhNQ2FYTXhDekFKQmdOVkJBZ01BbVJrTVFzd0NRWURWUVFIREFKa1pERUxNQWtHQTFVRUNnd0NaR1F4Q3pBSkJnTlZCQXNNQW1Sa01Rc3dDUVlEVlFRRERBSmtaREVSTUE4R0NTcUdTSWIzRFFFSkFSWUNaR1F3SGhjTk1qQXdOVEkxTVRZd05UQXdXaGNOTXpBd05USXpNVFl3TlRBd1dqQmhNUXN3Q1FZRFZRUUdFd0pwY3pFTE1Ba0dBMVVFQ0F3Q1pHUXhDekFKQmdOVkJBY01BbVJrTVFzd0NRWURWUVFLREFKa1pERUxNQWtHQTFVRUN3d0NaR1F4Q3pBSkJnTlZCQU1NQW1Sa01SRXdEd1lKS29aSWh2Y05BUWtCRmdKa1pEQ0NBU0l3RFFZSktvWklodmNOQVFFQkJRQURnZ0VQQURDQ0FRb0NnZ0VCQU1MNjNDWGtCYitsdnJKS2ZkZllCSExEWWZ1YUM2ZXhDU3FBU1VBb3NKV1dyZnlEaURNVWJtZnMwNlBMS3l2N044ZWZEaHphNzRvdjBFUUpOUmhNTmFDRStBMGNlcTZaWG1tTXN3VVlGZExBeThLMlZNejVtcm9CRlg4c2o1UFdWcjZyREoyY2tCYUZLV0JCOE5GbWlLN01UV1NJRjluOE0xMDcvOWEwUVVSQ3ZUaFVZdStzZ3V6YnNMT0RGdFhVeEc1cnRUVktCVmNQWnZFZlJreTJUa3Q0QXlTRlNta082S2Y0c0JkN01DNG1LV1ptN0s4azdIclpZejJ1c1NwYnJFdFlHdHI2TW1OOWhjaSsvSVREUEUyOTFERmt6SWNEQ0Y0OTN2LzNUKzdYc25tUWFqaDZrdUkrYmpJYUFDZm84Tit0d0VvSmYvTjFQbXBoQVFkRWlDMENBd0VBQWFOVE1GRXdIUVlEVlIwT0JCWUVGTnZtU3ByUVEySFVVdFB4czZVT3V4cTlsS0twTUI4R0ExVWRJd1FZTUJhQUZOdm1TcHJRUTJIVVV0UHhzNlVPdXhxOWxLS3BNQThHQTFVZEV3RUIvd1FGTUFNQkFmOHdEUVlKS29aSWh2Y05BUUVMQlFBRGdnRUJBSkVXeG54dFFWNUlxUFZScjJTTVdOTnhjSjdBL3d5ZXQzOWw1VmhIamJyUUd5bms1V1M4MHBzbi9yaUxVZkl2dHpZTVdDMElSMHBJTVF1TURGNXNOY0twNEQ4WG5yZCtCbC80L0l5L2lUT29IbHcrc1BrS3YrTkwyWFIzaU84YlNEd2p0anZkNkw1TmtVdXpzUm9Ta1FDRzJmSEFTcXFnRm95VjlMZFJzUWExdzlaR2VidEVXTHVHc3JKdFI3Z2FGRUNxSm5EYmIwYVBVTWl4bXBNSElEOGt0MTU0VHJMaFZGbU1FcUdHQzFHdlpWbFE5T2YzR1A5eTdYNHZEcEhzaGRsV290T25ZS0hhZXUyZDVjUlZGSGhFYnJzbGtJU2doL1RSdXlsN1ZJcG5qT1lVd01CcENpVkg2TTJseURJNlVSM0ZiejRwVlZBeEdYblZoQkV4akJFPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tLS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURRRENDQWlpZ0F3SUJBZ0lJZGd5RFlHbkRoVkl3RFFZSktvWklodmNOQVFFTEJRQXdQakVTTUJBR0ExVUUKQ3hNSmIzQmxibk5vYVdaME1TZ3dKZ1lEVlFRREV4OXJkV0psTFdGd2FYTmxjblpsY2kxc2IyTmhiR2h2YzNRdApjMmxuYm1WeU1CNFhEVEl3TURVeU5qRTFNRFV5TjFvWERUTXdNRFV5TkRFMU1EVXlOMW93UGpFU01CQUdBMVVFCkN4TUpiM0JsYm5Ob2FXWjBNU2d3SmdZRFZRUURFeDlyZFdKbExXRndhWE5sY25abGNpMXNiMk5oYkdodmMzUXQKYzJsbmJtVnlNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQTNzdnNGbVlYd3pNMApIcDNDY1p3Z2RBSnJWMXRocDNxMzcyMFV1RkdIUitvSXhhNkkwQXBVcnlxalhnY3JBSXRSTzBOZlkzMGUvaXdICmZ6bXI0MTZ4YVg0S3JXNUl1SkIyZkQ1OVFUQVBPRUJJSWhrRkV2Z0cxMitQaGlOVmhsbGt2V0oyTGVIL1V3cjMKN1J3eHlERFFrQ3ZDR1o5Ky83djlmSmVVVnJhSW5sL1dYTGp0c2orRXowRm1Ucks3bXpKaVFQcGs1RjNmUHpyMwpHQlh6REZQSU9JeHR6TTdsSng2dFNlaG5hZlByTHMxb3lXTm1HdXRNK0trYzE1dUROWUhybUo4Q1ZqZm83a2Q0CnQ3SXR6NUN5TGxwRFdwVEdEdmVoOXdMRnpRVjN2NHA0amJNZlZVTVJiV1BFRHp4cWphVGIwalVvVUFBRlBicTkKUTMxNG1zWlF3d0lEQVFBQm8wSXdRREFPQmdOVkhROEJBZjhFQkFNQ0FxUXdEd1lEVlIwVEFRSC9CQVV3QXdFQgovekFkQmdOVkhRNEVGZ1FVWkNGdmF3QWdWSnZtcGhIT20vbEdlZXdjZzI4d0RRWUpLb1pJaHZjTkFRRUxCUUFECmdnRUJBTmNxY0ZaMEJnZmM0TnlhclMwNXRWeGkzbzU2Qi9KMmpwNnVRVUtvZlJyMGhraUhYVDd0QnZubEppUSsKeSt5OHBEb29pVHpJOHMrMFd1dEN3S1d5dHdXNnNIczd4Zk9MUEtncFZXZWgyekRyZDhHY29hbXhlMzRmalRwdgppeUxDeWN2TkIzb1FpWkcyWkU1SDRXTmMwUmZXU0ZHUDcweG1JKzVVY2RvZTR6R3lEQU43bm5CYmZtU2xIK0kwCkkzWC9FaS9lRnZka081eG9kYVlUbFNWV2hTOHdib3FDVjF0SEN0YmIyQjYvY1VWTTRXUHg2Nm81Q05PMmYzMHoKcWJMU3F4ZzhYUXhjYW5FTGsxaGhNTjVlZDZkanc2YVFZUWxGSXBlcklPVVVTNDNJMGpRZTBYdmsrWWJhZGJscgpUVVR2VVV6RTdTNW5MUGRKK0hYZ2hZdFBHR1k9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0KLS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURURENDQWpTZ0F3SUJBZ0lJSXFpUUtKdUhNY293RFFZSktvWklodmNOQVFFTEJRQXdSREVTTUJBR0ExVUUKQ3hNSmIzQmxibk5vYVdaME1TNHdMQVlEVlFRREV5VnJkV0psTFdGd2FYTmxjblpsY2kxelpYSjJhV05sTFc1bApkSGR2Y21zdGMybG5ibVZ5TUI0WERUSXdNRFV5TmpFMU1EVXlOMW9YRFRNd01EVXlOREUxTURVeU4xb3dSREVTCk1CQUdBMVVFQ3hNSmIzQmxibk5vYVdaME1TNHdMQVlEVlFRREV5VnJkV0psTFdGd2FYTmxjblpsY2kxelpYSjIKYVdObExXNWxkSGR2Y21zdGMybG5ibVZ5TUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQwpBUUVBM3FDYjEycHFMcEt3M1dTVU1MclhTdFcza2ltYllVWGk1R1AyeXJhZjFYOE5pMlZwZnZHcUxtaWFpNFhVCjF2WkJhdEhYSE1SS1VteDczUTA3ZWVKUnQzNTdSUU00ZW9tbVhqQUloQjhPRU1JV0gzTTN1dTFnVlUvZkpNb2wKTlozY3JvN1BjUGFtVmRzdzIwYjhSTEIwcW55NW9FSWdLU2F5YkdGWlF6V2t4WitaR2NlNWxLbmZ0OGZpci9Fegp3NkpOK0lxbW1YOTFtVXR6T25CaFZNNFB5RzFHbnUyZ1g1SXNUOWFHbGdwOWlTSENJMnZyMEtORXZVZFJaY0k2Ck9SM0ZQUSs1WnQ5c2hPOHhYZlJpMFF3SmN2bG5xdCtDS2N6RFcwelUwc1BhVzJ5VlJzUmIyQjZFVHExMGxpQXoKN3NKRzdFc3U4TW5JVnQ1dFpqMmoxSkpUMVFJREFRQUJvMEl3UURBT0JnTlZIUThCQWY4RUJBTUNBcVF3RHdZRApWUjBUQVFIL0JBVXdBd0VCL3pBZEJnTlZIUTRFRmdRVXZ4RGJ5NnJaNWRkczBTOEtDR3dNcVRhWitFb3dEUVlKCktvWklodmNOQVFFTEJRQURnZ0VCQUtIdmxZR2N2NTFyQW1QalR2dWFIZUpaQjhoRmlsTkdHdnIxK1QxajZIY2YKcndPYzYxSCtGUjJtenNqRFBMSWhBSjV1WWlidUZEQVB3eldDUXpJT0lxeDZhanhaZlhaSTNBcDdibDB1TlBtUwpEeW1XdnJVemx2MGZoMHVwUHE1RGFiNXFPMnRPTWJxTjNjVDQxNXU0YldWaThnVmZ2Mk9SWFUxVEhpWkliVUVJCkVNbE1vRStSRU5WQzZlU0U3Rmd3S3lUQ2M2cmZ4ZGtyMU8yOEhYOFNOOGhydW9WYXdJajN6Z1k5ZWVlZUpxdE0KSmdPVVRwWThKSWIzRzcwZjRId01xQ3hrcERqUWROVTVMYTR1OUlBN0NlQldiemdrTVBiMmI4bmRMUngzQ0Z2MAovcmF1U21GY2l0MnVid1BYbWIzalBFcVhuWXh1ckp0YWtvVTAvQnRraGdrPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCi0tLS0tQkVHSU4gQ0VSVElGSUNBVEUtLS0tLQpNSUlETWpDQ0FocWdBd0lCQWdJSUhTRXB2RFlvcHBnd0RRWUpLb1pJaHZjTkFRRUxCUUF3TnpFU01CQUdBMVVFCkN4TUpiM0JsYm5Ob2FXWjBNU0V3SHdZRFZRUURFeGhyZFdKbExXRndhWE5sY25abGNpMXNZaTF6YVdkdVpYSXcKSGhjTk1qQXdOVEkyTVRVd05USTNXaGNOTXpBd05USTBNVFV3TlRJM1dqQTNNUkl3RUFZRFZRUUxFd2x2Y0dWdQpjMmhwWm5ReElUQWZCZ05WQkFNVEdHdDFZbVV0WVhCcGMyVnlkbVZ5TFd4aUxYTnBaMjVsY2pDQ0FTSXdEUVlKCktvWklodmNOQVFFQkJRQURnZ0VQQURDQ0FRb0NnZ0VCQU1nREt3Ni9ONkY2UXBWZEJRdXpHSy95RFhkUGhtRXgKNWNiK2h6ZEw2RDFPNWZVRHBjZitHT1FEdjB2dUZrNXcveDlVRmxoRldRTkNHOEtQbWxvSFhkMHIxQllWYWc0VgpYcGFYYklGUHJmMzVGeFhzZk1HS0JUQllaRXNiKzJRTWFNci9Ec0g0cTQ5Z294NDJ5U24rV2prbnNCWnZ4Ukd0CndBT2t0b1VWUVloaWxrWDg5dHl1VjJlTU1CbmVEZmdyazZ1WEt4YjIzTWg2WEZvNHAvMW9MV1VrVmxQQ3hlM2MKNWI2VWNON29NdTI3M016enJnZzVDOXdmbGVycWFoamxBRXFHY0JMQ3NrT0NBbGhMWGRvMW5FZ09lM0ZMUHJFeApiNW93UTJkNnQvbEFlellhbU5lTUl4bTN1OWx6djJXVytLUm4rRmcxci9UZWt4Qi9ZUUVkNHhjQ0F3RUFBYU5DCk1FQXdEZ1lEVlIwUEFRSC9CQVFEQWdLa01BOEdBMVVkRXdFQi93UUZNQU1CQWY4d0hRWURWUjBPQkJZRUZHN0QKOXc2UEloYUxEYTRIaU05MVNMdUw3bGlUTUEwR0NTcUdTSWIzRFFFQkN3VUFBNElCQVFBV2JZeENVUzF2SStEUQpvK1JJZmxRSU5KZFcwZkF5WEFtVUdPZGdiVW8zL2R5UldnY0RDcFdCbVhMWlhkaVJwd3R5MWV2NTNzM1VrM0ZrCldoRGJSZVFPb29nSW1RMmVZdnZxWjFlbHhGeHdYdHZtVGw1cGZ1aFdMSDI4UC8rZnZFQW04T1dpTmtSUmpsNWYKMkhLUEU1MzBMQmtsSGx0OE9naUh0bS9XOGRUcTk3Y3lXaVlsdlREeFRkREFPRTVLZERrNXpxbXRBWUpWRlhnWApyajR0UUthTkphRWdyK2JkamhKbUIyUjU4b3RxMExXeS9uNnBjTWNmSFk2Q2k5dWhwT0Y3R3RWYjdTcUNGd0tPClMvdlVORFBsK2QyRjdLMDQweUVRUExQZHI5cEF5QlNWSkZyL2lnRkJycHJJVEdkYUZ1dFhQUHYybjBtTHBvS0kKVTlIV2FJcTgKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
      server: "https://api.test-infra-cluster.redhat:6443"
    name: test-infra-cluster
contexts: 
  - 
    context: 
      cluster: test-infra-cluster
      user: admin
    name: admin
current-context: admin
kind: Config
preferences: {}
users: 
  - 
    name: admin
    user: 
      client-certificate-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURaekNDQWsrZ0F3SUJBZ0lJUnljNWdMZE1sMjB3RFFZSktvWklodmNOQVFFTEJRQXdOakVTTUJBR0ExVUUKQ3hNSmIzQmxibk5vYVdaME1TQXdIZ1lEVlFRREV4ZGhaRzFwYmkxcmRXSmxZMjl1Wm1sbkxYTnBaMjVsY2pBZQpGdzB5TURBMU1qWXhOVEExTWpaYUZ3MHpNREExTWpReE5UQTFNamRhTURBeEZ6QVZCZ05WQkFvVERuTjVjM1JsCmJUcHRZWE4wWlhKek1SVXdFd1lEVlFRREV3eHplWE4wWlcwNllXUnRhVzR3Z2dFaU1BMEdDU3FHU0liM0RRRUIKQVFVQUE0SUJEd0F3Z2dFS0FvSUJBUURXR0h2ak13Tzh6SGRPTGtJcmVYMFYvdXptMjZ6dlJYQWJ2K3R6YU42dgpXbFByTFZUSkRwVHJwUVpmbnI4ZlBXQ0xSazRmRlNTZlZlR2EvRi80MzNUNDdTWGpQcHhYWEV3Y1VhYi8yRjhhCnlmdlphVGRXRHBNd3NFNzYzeXBEODhvalNQQVlRdDRqSGhRcEV3TExqVVR6VXdkTGJjdnREVEQyWmtwTE1oZFcKMmhrNGU1eVZXSzhrUDRsSUxYbnJMbFJxR2VUNGRtbmN6NHdMcTZSWXgrR0hUaFhNUzlTZEhkZ2V1dHlVZGlVMQpPNFRIOGNEZklrRnQ5Y3lkSmNZajdmREh6WWRQemlVdmtrcFJ0eHNuamlFTzVwM0h0cDAyL3kvUmlpN2t5TjByCndlNGM2Y29LZzRXZzYrUzVzaEpZV3d0ZGsrTmkxYU1BTnVIOExMbGhaUnlEQWdNQkFBR2pmekI5TUE0R0ExVWQKRHdFQi93UUVBd0lGb0RBZEJnTlZIU1VFRmpBVUJnZ3JCZ0VGQlFjREFRWUlLd1lCQlFVSEF3SXdEQVlEVlIwVApBUUgvQkFJd0FEQWRCZ05WSFE0RUZnUVVqWWlDdGo4blg2dHIzTHFCTURmSHU0aDdERUF3SHdZRFZSMGpCQmd3CkZvQVUxYkQ2VUUwN3cxUFBHMzZKWjJSZHJLL2lQUEl3RFFZSktvWklodmNOQVFFTEJRQURnZ0VCQUhHUWpSRUoKU3NsbFM1YkdHYkxmNGFGdVIxbDdiVGh0UTBLTFl3bkowQmtKNDBnVVZUcHdmMEFpcFRVanU1UTkxNzZaNEJWUQpleFB4U1BuRngwSGxoQXhUYmJQaWQ2Y25leVlId2I1Vi95R3I2M0p5cWcxR2pScVJ5djBLR3k0bHRnWGRKNmQyCktZaVFlT2pYN3B2SXVSMjE2T2ZNMWdnYUdOUkFzSDZwYjlEMU5ENmoyRnVlUkcrVFN1MUxzRlJXWUZsVHk0aDMKc014RElWblY3L0R1T1M5QTJCVm80alpYL0pmbzFId2ExL3BaVkJ1YVU5UHJlemNOblBkZ09TRmxVc0tBbUNjMgpCb3dES0RPbkZZMm9IZTc5ZjF6NlU1a1Buby80ZGlBT0FIQ2dPY0lPbk9ubUNOR1d3R0tLU2tGc05lYUtnVVltCjVYVkFqNVd3amNTdmhqMD0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
      client-key-data: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb3dJQkFBS0NBUUVBMWhoNzR6TUR2TXgzVGk1Q0szbDlGZjdzNXR1czcwVndHNy9yYzJqZXIxcFQ2eTFVCnlRNlU2NlVHWDU2L0h6MWdpMFpPSHhVa24xWGhtdnhmK045MCtPMGw0ejZjVjF4TUhGR20vOWhmR3NuNzJXazMKVmc2VE1MQk8rdDhxUS9QS0kwandHRUxlSXg0VUtSTUN5NDFFODFNSFMyM0w3UTB3OW1aS1N6SVhWdG9aT0h1YwpsVml2SkQrSlNDMTU2eTVVYWhuaytIWnAzTStNQzZ1a1dNZmhoMDRWekV2VW5SM1lIcnJjbEhZbE5UdUV4L0hBCjN5SkJiZlhNblNYR0krM3d4ODJIVDg0bEw1SktVYmNiSjQ0aER1YWR4N2FkTnY4djBZb3U1TWpkSzhIdUhPbksKQ29PRm9Pdmt1YklTV0ZzTFhaUGpZdFdqQURiaC9DeTVZV1VjZ3dJREFRQUJBb0lCQUFRRHo0YnlOUGE4YXR4WApkN3d5K2dxSWprN0IvZHM2MVNCZ0YvMUJFVFArb0tZL1ltQ20ybG9VN1NxcjRtK21pZ0h5bnBKc3BoUXEyeUU1CjdGN1JhZk1sRjFuTW1jZjFuaVBGMERqcUNOYUt4U05ObXRFTlV1dE4weDFYUkFha01yMDRwKy84aVFmbGo0RTUKcndxOEtuZlpyY0JYWGNTalE3RExPRWR5dUFkVDVPN0lqRjV3dENNNUw0Tlc5K0dkOFg4VnpxaTNndWlkdzUwRQphSjlndjR2Nnp3RDAydFNSZHNoUlBBbU9oVzBCVDl5bkZCNEFrdlg3b1psVHR2YU5QMVF3S0RkQXd0cmF5K2ZXCjRMcmY0ZldEWDIzNERqS01sOUN3cjhnTmNja0FoejVjN09HQ0dzeU9xc3JBTmllNFB6VWtmc0xhM1ZQQ0FsTHoKV2s1dEs0a0NnWUVBNVBMM3hzMzB2Y1NPcjhjV2R4eElCTE1pTzNnRldBb05qeEJkQ2ExNEZCYjlGaVMxK1NNSgpQTkZNTU1rT2VUMGZPamx1dVpyTGVueDIrUFZNVnhrb29kYXNmcGxtNGI3RHk1RkZzcGhUajlOOXJJdmIvWWpJCk51SDdKNkgzVEZJTzhDVE5GK1FwSmIza2pNZUhRcWM0S0cwSFBCeTlpLytvUExkOGdxZ1JCczhDZ1lFQTcyUSsKT3prYlZaUzlad25YcEZPZmQwa2hZQkkzUFlvaElockJjOEw3Q2xHZlVwdFRQdlMxSWVqUjlMNHFiM1BsSHB6bAptaG5wQXRGTFFYQXJCdFk4UmErZkxlbW1Xdzc2ZndrdlFraXlLblUvZ0s1UlVkdmdQbkxrSCtVSmJJTmwzcjgyCkRzQzB0bWJJL1QrblpObEJHdVFkQm1US3hkZ1R1MXRBZ3dIcy9BMENnWUVBMlZISDMrMmZZb0l3N3FrTHFnUXUKV0VleE5zRzJVTnM2QTVLRXZhcnJVQ2FDRllMRE9Ma0pDN0dmb0s4NERkejJ4MDI4ekhFaXRDRnd6T0FLbHFKSwo3MVBXYUZVMFV4UEF4bm9lcm1mbzZaeldyZklUMzVUMmR5SUtSSlI1S1BpN05UZTVkZlFkR3JZbE8zd3A2QnJTCk00MUtVTVQzSnV5RnhSeG1FNTkwaWdFQ2dZQjU5bHRTTnVUN00vMU8rbyszczdiaHdndFQ4OVBhOFgyeDcybXgKdlp2Q2hSVWpzK2kwZ1YycStmL0ZyZ0RXcVhnSW9hekVWd0VFbzNhd3p5SE1xT2NxSmJCMlpyeVBWZEUvV1lHUApScFFtMTNkVDZ2dVpOZWxJUjZaN3JXZWd0a3ozTC9tdGlIWkpHNUs0bTI2QURjT0NuTWRBMDZjUEp1ZmVvejM1CndNaHBIUUtCZ0E1WVpYeWE1azUrcGFRd1RaWTFKb2ZsZ2FjNlRRZ2IzOTJDWStZR1RjNHdDVVlNNWVVazg0TCsKRFM4bTYwc0liL0g3UVBOa2RBTXZobFNPMk9McFVEbGtWUnhrUXFvckxWRDRRQUc0SFBMU1VLVzI0eFNqeXpZNApncXQ4b3VmcEFIV0pMLzZpZEthUzhvVjVlc2NyV0Y5clJ6bWliUGx4V1JIZ3ZzdVJJSEhjCi0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg==
`

const (
	fakeIgnitionConfigOverride      = `{"ignition": {"version": "3.1.0"}, "storage": {"files": [{"path": "/tmp/example", "contents": {"source": "data:text/plain;base64,aGVscGltdHJhcHBlZGluYXN3YWdnZXJzcGVj"}}]}}`
	badIgnitionConfigOverride       = `bad ignition config`
	adminKubeConfigStringTemplate   = "%s-admin-kubeconfig"
	BMH_HARDWARE_DETAILS_ANNOTATION = "inspect.metal3.io/hardwaredetails"
	MACHINE_ROLE                    = "machine.openshift.io/cluster-api-machine-role"
	MACHINE_TYPE                    = "machine.openshift.io/cluster-api-machine-type"
)

var (
	imageSetsData = map[string]string{
		"openshift-v4.7.0": "quay.io/openshift-release-dev/ocp-release:4.7.2-x86_64",
		"openshift-v4.8.0": "quay.io/openshift-release-dev/ocp-release:4.8.0-fc.0-x86_64",
	}
)

func deployLocalObjectSecretIfNeeded(ctx context.Context, client k8sclient.Client) *corev1.LocalObjectReference {
	err := client.Get(
		ctx,
		types.NamespacedName{Namespace: Options.Namespace, Name: "pull-secret"},
		&corev1.Secret{},
	)
	if apierrors.IsNotFound(err) {
		deployPullSecretResource(ctx, kubeClient, "pull-secret", pullSecret)
	} else {
		Expect(err).To(BeNil())
	}
	return &corev1.LocalObjectReference{
		Name: "pull-secret",
	}
}

func deployPullSecretResource(ctx context.Context, client k8sclient.Client, name, secret string) {
	data := map[string]string{corev1.DockerConfigJsonKey: secret}
	s := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: Options.Namespace,
			Name:      name,
		},
		StringData: data,
		Type:       corev1.SecretTypeDockerConfigJson,
	}
	Expect(client.Create(ctx, s)).To(BeNil())
}

func deployClusterDeploymentCRD(ctx context.Context, client k8sclient.Client, spec *hivev1.ClusterDeploymentSpec) {
	deployClusterImageSetCRD(ctx, client, spec.Provisioning.ImageSetRef)
	err := client.Create(ctx, &hivev1.ClusterDeployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterDeployment",
			APIVersion: getAPIVersion(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: Options.Namespace,
			Name:      spec.ClusterName,
		},
		Spec: *spec,
	})
	Expect(err).To(BeNil())
}

func deployBMHCRD(ctx context.Context, client k8sclient.Client, name string, spec *bmhv1alpha1.BareMetalHostSpec) {
	err := client.Create(ctx, &bmhv1alpha1.BareMetalHost{
		TypeMeta: metav1.TypeMeta{
			Kind:       "BareMetalHost",
			APIVersion: "metal3.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: Options.Namespace,
			Name:      name,
		},
		Spec: *spec,
	})
	Expect(err).To(BeNil())
}

func addAnnotationToClusterDeployment(ctx context.Context, client k8sclient.Client, key types.NamespacedName, annotationKey string, annotationValue string) {
	Eventually(func() error {
		clusterDeploymentCRD := getClusterDeploymentCRD(ctx, client, key)
		clusterDeploymentCRD.SetAnnotations(map[string]string{annotationKey: annotationValue})
		return kubeClient.Update(ctx, clusterDeploymentCRD)
	}, "30s", "10s").Should(BeNil())
}

func deployClusterImageSetCRD(ctx context.Context, client k8sclient.Client, imageSetRef *hivev1.ClusterImageSetReference) {
	err := client.Create(ctx, &hivev1.ClusterImageSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterImageSet",
			APIVersion: getAPIVersion(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: Options.Namespace,
			Name:      imageSetRef.Name,
		},
		Spec: hivev1.ClusterImageSetSpec{
			ReleaseImage: imageSetsData[imageSetRef.Name],
		},
	})
	Expect(err).To(BeNil())
}

func deployInfraEnvCRD(ctx context.Context, client k8sclient.Client, name string, spec *v1beta1.InfraEnvSpec) {
	err := client.Create(ctx, &v1beta1.InfraEnv{
		TypeMeta: metav1.TypeMeta{
			Kind:       "InfraEnv",
			APIVersion: getAPIVersion(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: Options.Namespace,
			Name:      name,
		},
		Spec: *spec,
	})

	Expect(err).To(BeNil())
}

func deployNMStateConfigCRD(ctx context.Context, client k8sclient.Client, name string, NMStateLabelName string, NMStateLabelValue string, spec *v1beta1.NMStateConfigSpec) {
	err := client.Create(ctx, &v1beta1.NMStateConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       "NMStateConfig",
			APIVersion: getAPIVersion(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: Options.Namespace,
			Name:      name,
			Labels:    map[string]string{NMStateLabelName: NMStateLabelValue},
		},
		Spec: *spec,
	})
	Expect(err).To(BeNil())
}

func getAPIVersion() string {
	return fmt.Sprintf("%s/%s", v1beta1.GroupVersion.Group, v1beta1.GroupVersion.Version)
}

func getClusterFromDB(
	ctx context.Context, client k8sclient.Client, db *gorm.DB, key types.NamespacedName, timeout int) *common.Cluster {

	var err error
	cluster := &common.Cluster{}
	start := time.Now()
	for time.Duration(timeout)*time.Second > time.Since(start) {
		cluster, err = common.GetClusterFromDBWhere(db, common.UseEagerLoading, common.SkipDeletedRecords, "kube_key_name = ? and kube_key_namespace = ?", key.Name, key.Namespace)
		if err == nil {
			return cluster
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			Expect(err).To(BeNil())
		}
		getClusterDeploymentCRD(ctx, client, key)
		time.Sleep(time.Second)
	}
	Expect(err).To(BeNil())
	return cluster
}

func getClusterDeploymentCRD(ctx context.Context, client k8sclient.Client, key types.NamespacedName) *hivev1.ClusterDeployment {
	cluster := &hivev1.ClusterDeployment{}
	err := client.Get(ctx, key, cluster)
	Expect(err).To(BeNil())
	return cluster
}

func getInfraEnvCRD(ctx context.Context, client k8sclient.Client, key types.NamespacedName) *v1beta1.InfraEnv {
	infraEnv := &v1beta1.InfraEnv{}
	err := client.Get(ctx, key, infraEnv)
	Expect(err).To(BeNil())
	return infraEnv
}

func getAgentCRD(ctx context.Context, client k8sclient.Client, key types.NamespacedName) *v1beta1.Agent {
	agent := &v1beta1.Agent{}
	err := client.Get(ctx, key, agent)
	Expect(err).To(BeNil())
	return agent
}

func getBmhCRD(ctx context.Context, client k8sclient.Client, key types.NamespacedName) *bmhv1alpha1.BareMetalHost {
	bmh := &bmhv1alpha1.BareMetalHost{}
	err := client.Get(ctx, key, bmh)
	Expect(err).To(BeNil())
	return bmh
}

func getSecret(ctx context.Context, client k8sclient.Client, key types.NamespacedName) *corev1.Secret {
	secret := &corev1.Secret{}
	err := client.Get(ctx, key, secret)
	Expect(err).To(BeNil())
	return secret
}

// configureLoclAgentClient reassigns the global agentBMClient variable to a client instance using local token auth
func configureLocalAgentClient(clusterID string) {
	if Options.AuthType != auth.TypeLocal {
		Fail(fmt.Sprintf("Agent client shouldn't be configured for local auth when auth type is %s", Options.AuthType))
	}

	key := types.NamespacedName{
		Namespace: Options.Namespace,
		Name:      "assisted-installer-local-auth-key",
	}
	secret := getSecret(context.Background(), kubeClient, key)
	privKeyPEM := secret.Data["ec-private-key.pem"]
	tok, err := gencrypto.LocalJWTForKey(clusterID, string(privKeyPEM))
	Expect(err).To(BeNil())

	agentBMClient = client.New(clientcfg(auth.AgentAuthHeaderWriter(tok)))
}

func checkAgentCondition(ctx context.Context, hostId string, conditionType conditionsv1.ConditionType, reason string) {
	hostkey := types.NamespacedName{
		Namespace: Options.Namespace,
		Name:      hostId,
	}
	Eventually(func() string {
		return conditionsv1.FindStatusCondition(getAgentCRD(ctx, kubeClient, hostkey).Status.Conditions, conditionType).Reason
	}, "30s", "10s").Should(Equal(reason))
}

func checkClusterCondition(ctx context.Context, key types.NamespacedName, conditionType hivev1.ClusterDeploymentConditionType, reason string) {
	Eventually(func() string {
		condition := controllers.FindStatusCondition(getClusterDeploymentCRD(ctx, kubeClient, key).Status.Conditions, conditionType)
		if condition != nil {
			return condition.Reason
		}
		return ""
	}, "2m", "2s").Should(Equal(reason))
}

func checkInfraEnvCondition(ctx context.Context, key types.NamespacedName, conditionType conditionsv1.ConditionType, message string) {
	Eventually(func() string {
		condition := conditionsv1.FindStatusCondition(getInfraEnvCRD(ctx, kubeClient, key).Status.Conditions, conditionType)
		if condition != nil {
			return condition.Message
		}
		return ""
	}, "2m", "1s").Should(Equal(message))
}

func getDefaultClusterDeploymentSpec(secretRef *corev1.LocalObjectReference) *hivev1.ClusterDeploymentSpec {
	return &hivev1.ClusterDeploymentSpec{
		ClusterName: "test-cluster",
		BaseDomain:  "hive.example.com",
		Provisioning: &hivev1.Provisioning{
			InstallConfigSecretRef: &corev1.LocalObjectReference{Name: "cluster-install-config"},
			ImageSetRef:            &hivev1.ClusterImageSetReference{Name: "openshift-v4.8.0"},
			InstallStrategy: &hivev1.InstallStrategy{
				Agent: &agentv1.InstallStrategy{
					Networking: agentv1.Networking{
						MachineNetwork: []agentv1.MachineNetworkEntry{},
						ClusterNetwork: []agentv1.ClusterNetworkEntry{{
							CIDR:       "10.128.0.0/14",
							HostPrefix: 23,
						}},
						ServiceNetwork: []string{"172.30.0.0/16"},
					},
					SSHPublicKey: sshPublicKey,
					ProvisionRequirements: agentv1.ProvisionRequirements{
						ControlPlaneAgents: 3,
						WorkerAgents:       0,
					},
				},
			},
		},
		Platform: hivev1.Platform{
			AgentBareMetal: &agentv1.BareMetalPlatform{
				APIVIP:     "1.2.3.8",
				IngressVIP: "1.2.3.9",
			},
		},
		PullSecretRef: secretRef,
	}
}

func getDefaultClusterDeploymentSNOSpec(secretRef *corev1.LocalObjectReference) *hivev1.ClusterDeploymentSpec {
	return &hivev1.ClusterDeploymentSpec{
		ClusterName: "test-cluster-sno",
		BaseDomain:  "hive.example.com",
		Provisioning: &hivev1.Provisioning{
			InstallConfigSecretRef: &corev1.LocalObjectReference{Name: "cluster-install-config"},
			ImageSetRef:            &hivev1.ClusterImageSetReference{Name: "openshift-v4.8.0"},
			InstallStrategy: &hivev1.InstallStrategy{
				Agent: &agentv1.InstallStrategy{
					Networking: agentv1.Networking{
						MachineNetwork: []agentv1.MachineNetworkEntry{{CIDR: "1.2.3.0/24"}},
						ClusterNetwork: []agentv1.ClusterNetworkEntry{{
							CIDR:       "10.128.0.0/14",
							HostPrefix: 23,
						}},
						ServiceNetwork: []string{"172.30.0.0/16"},
					},
					SSHPublicKey: sshPublicKey,
					ProvisionRequirements: agentv1.ProvisionRequirements{
						ControlPlaneAgents: 1,
						WorkerAgents:       0,
					},
				},
			},
		},
		Platform: hivev1.Platform{
			AgentBareMetal: &agentv1.BareMetalPlatform{},
		},
		PullSecretRef: secretRef,
	}
}

func getDefaultInfraEnvSpec(secretRef *corev1.LocalObjectReference,
	clusterDeployment *hivev1.ClusterDeploymentSpec) *v1beta1.InfraEnvSpec {
	return &v1beta1.InfraEnvSpec{
		ClusterRef: &v1beta1.ClusterReference{
			Name:      clusterDeployment.ClusterName,
			Namespace: Options.Namespace,
		},
		PullSecretRef:    secretRef,
		SSHAuthorizedKey: sshPublicKey,
	}
}

func getDefaultNMStateConfigSpec(nicPrimary, nicSecondary, macPrimary, macSecondary, networkYaml string) *v1beta1.NMStateConfigSpec {
	return &v1beta1.NMStateConfigSpec{
		Interfaces: []*v1beta1.Interface{
			{MacAddress: macPrimary, Name: nicPrimary},
			{MacAddress: macSecondary, Name: nicSecondary},
		},
		NetConfig: v1beta1.NetConfig{Raw: []byte(networkYaml)},
	}
}

func getAgentMac(agent *v1beta1.Agent) string {

	for _, agentInterface := range agent.Status.Inventory.Interfaces {
		if agentInterface.MacAddress != "" {
			return agentInterface.MacAddress
		}
	}
	return ""
}

func cleanUP(ctx context.Context, client k8sclient.Client) {
	Expect(client.DeleteAllOf(ctx, &hivev1.ClusterDeployment{}, k8sclient.InNamespace(Options.Namespace))).To(BeNil())
	Expect(client.DeleteAllOf(ctx, &hivev1.ClusterImageSet{}, k8sclient.InNamespace(Options.Namespace))).To(BeNil())
	Expect(client.DeleteAllOf(ctx, &v1beta1.InfraEnv{}, k8sclient.InNamespace(Options.Namespace))).To(BeNil())
	Expect(client.DeleteAllOf(ctx, &v1beta1.NMStateConfig{}, k8sclient.InNamespace(Options.Namespace))).To(BeNil())
	Expect(client.DeleteAllOf(ctx, &v1beta1.Agent{}, k8sclient.InNamespace(Options.Namespace))).To(BeNil())
	Expect(client.DeleteAllOf(ctx, &bmhv1alpha1.BareMetalHost{}, k8sclient.InNamespace(Options.Namespace))).To(BeNil())
	// Expect(client.DeleteAllOf(ctx, &corev1.Secret{}, k8sclient.InNamespace(Options.Namespace))).To(BeNil())
}

func setupNewHost(ctx context.Context, hostname string, clusterID strfmt.UUID) *models.Host {
	host := registerNode(ctx, clusterID, hostname)
	generateHWPostStepReply(ctx, host, validHwInfo, hostname)
	generateFAPostStepReply(ctx, host, validFreeAddresses)
	generateDiskSpeedResponses(ctx, sdbId, host)
	return host
}

var _ = Describe("[kube-api]cluster installation", func() {
	if !Options.EnableKubeAPI {
		return
	}

	ctx := context.Background()

	waitForReconcileTimeout := 30

	AfterEach(func() {
		cleanUP(ctx, kubeClient)
		clearDB()
	})

	/* It("deploy clusterDeployment with agents and wait for ready", func() {
		secretRef := deployLocalObjectSecretIfNeeded(ctx, kubeClient)
		spec := getDefaultClusterDeploymentSpec(secretRef)
		deployClusterDeploymentCRD(ctx, kubeClient, spec)
		key := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      spec.ClusterName,
		}
		cluster := getClusterFromDB(ctx, kubeClient, db, key, waitForReconcileTimeout)
		configureLocalAgentClient(cluster.ID.String())
		hosts := make([]*models.Host, 0)
		for i := 0; i < 3; i++ {
			hostname := fmt.Sprintf("h%d", i)
			host := setupNewHost(ctx, hostname, *cluster.ID)
			hosts = append(hosts, host)
		}
		generateFullMeshConnectivity(ctx, "1.2.3.10", hosts...)
		for _, host := range hosts {
			hostkey := types.NamespacedName{
				Namespace: Options.Namespace,
				Name:      host.ID.String(),
			}
			Eventually(func() error {
				agent := getAgentCRD(ctx, kubeClient, hostkey)
				agent.Spec.Approved = true
				return kubeClient.Update(ctx, agent)
			}, "30s", "10s").Should(BeNil())
		}
		checkClusterCondition(ctx, key, controllers.ClusterReadyForInstallationCondition, controllers.ClusterAlreadyInstallingReason)
	})

	It("deploy clusterDeployment with agent and update agent", func() {
		secretRef := deployLocalObjectSecretIfNeeded(ctx, kubeClient)
		spec := getDefaultClusterDeploymentSpec(secretRef)
		deployClusterDeploymentCRD(ctx, kubeClient, spec)
		key := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      spec.ClusterName,
		}
		cluster := getClusterFromDB(ctx, kubeClient, db, key, waitForReconcileTimeout)
		configureLocalAgentClient(cluster.ID.String())
		host := setupNewHost(ctx, "hostname1", *cluster.ID)
		key = types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      host.ID.String(),
		}
		Eventually(func() error {
			agent := getAgentCRD(ctx, kubeClient, key)
			agent.Spec.Hostname = "newhostname"
			agent.Spec.Approved = true
			agent.Spec.InstallationDiskID = sdb.ID
			return kubeClient.Update(ctx, agent)
		}, "30s", "10s").Should(BeNil())

		Eventually(func() string {
			h, err := common.GetHostFromDB(db, cluster.ID.String(), host.ID.String())
			Expect(err).To(BeNil())
			return h.RequestedHostname
		}, "2m", "10s").Should(Equal("newhostname"))
		Eventually(func() string {
			h, err := common.GetHostFromDB(db, cluster.ID.String(), host.ID.String())
			Expect(err).To(BeNil())
			return h.InstallationDiskID
		}, "2m", "10s").Should(Equal(sdb.ID))
		Eventually(func() bool {
			return conditionsv1.IsStatusConditionTrue(getAgentCRD(ctx, kubeClient, key).Status.Conditions, controllers.SpecSyncedCondition)
		}, "2m", "10s").Should(Equal(true))
		Eventually(func() string {
			return getAgentCRD(ctx, kubeClient, key).Status.Inventory.SystemVendor.Manufacturer
		}, "2m", "10s").Should(Equal(validHwInfo.SystemVendor.Manufacturer))
	})

	It("deploy clusterDeployment with agent,bmh and ignition config override", func() {
		secretRef := deployLocalObjectSecretIfNeeded(ctx, kubeClient)
		spec := getDefaultClusterDeploymentSpec(secretRef)
		deployClusterDeploymentCRD(ctx, kubeClient, spec)
		key := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      spec.ClusterName,
		}
		cluster := getClusterFromDB(ctx, kubeClient, db, key, waitForReconcileTimeout)
		configureLocalAgentClient(cluster.ID.String())
		host := setupNewHost(ctx, "hostname1", *cluster.ID)
		key = types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      host.ID.String(),
		}

		agent := getAgentCRD(ctx, kubeClient, key)
		bmhSpec := bmhv1alpha1.BareMetalHostSpec{BootMACAddress: getAgentMac(agent)}
		deployBMHCRD(ctx, kubeClient, host.ID.String(), &bmhSpec)

		Eventually(func() error {
			bmh := getBmhCRD(ctx, kubeClient, key)
			bmh.SetAnnotations(map[string]string{controllers.BMH_AGENT_IGNITION_CONFIG_OVERRIDES: fakeIgnitionConfigOverride})
			return kubeClient.Update(ctx, bmh)
		}, "30s", "10s").Should(BeNil())

		Eventually(func() bool {
			h, err := common.GetHostFromDB(db, cluster.ID.String(), host.ID.String())
			Expect(err).To(BeNil())
			agent := getAgentCRD(ctx, kubeClient, key)

			if agent.Spec.IgnitionConfigOverrides == "" {
				return false
			}

			if h.IgnitionConfigOverrides == "" {
				return false
			}
			return reflect.DeepEqual(h.IgnitionConfigOverrides, agent.Spec.IgnitionConfigOverrides)
		}, "2m", "10s").Should(Equal(true))

		By("Clean ignition config override ")
		h, err := common.GetHostFromDB(db, cluster.ID.String(), host.ID.String())
		Expect(err).To(BeNil())
		Expect(h.IgnitionConfigOverrides).NotTo(BeEmpty())

		Eventually(func() error {
			bmh := getBmhCRD(ctx, kubeClient, key)
			bmh.SetAnnotations(map[string]string{controllers.BMH_AGENT_IGNITION_CONFIG_OVERRIDES: ""})
			return kubeClient.Update(ctx, bmh)
		}, "30s", "10s").Should(BeNil())

		By("Verify agent ignition config override were cleaned")
		Eventually(func() string {
			agent := getAgentCRD(ctx, kubeClient, key)
			return agent.Spec.IgnitionConfigOverrides
		}, "30s", "10s").Should(BeEmpty())

		By("Verify host ignition config override were cleaned")
		Eventually(func() int {
			h, err := common.GetHostFromDB(db, cluster.ID.String(), host.ID.String())
			Expect(err).To(BeNil())

			return len(h.IgnitionConfigOverrides)
		}, "2m", "10s").Should(Equal(0))
	})

	It("deploy clusterDeployment with agent and invalid ignition config", func() {
		secretRef := deployLocalObjectSecretIfNeeded(ctx, kubeClient)
		spec := getDefaultClusterDeploymentSpec(secretRef)
		deployClusterDeploymentCRD(ctx, kubeClient, spec)
		key := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      spec.ClusterName,
		}
		cluster := getClusterFromDB(ctx, kubeClient, db, key, waitForReconcileTimeout)
		configureLocalAgentClient(cluster.ID.String())
		host := setupNewHost(ctx, "hostname1", *cluster.ID)
		key = types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      host.ID.String(),
		}

		h, err := common.GetHostFromDB(db, cluster.ID.String(), host.ID.String())
		Expect(err).To(BeNil())
		Expect(h.IgnitionConfigOverrides).To(BeEmpty())

		By("Invalid ignition config - invalid json")
		Eventually(func() error {
			agent := getAgentCRD(ctx, kubeClient, key)
			agent.Spec.IgnitionConfigOverrides = badIgnitionConfigOverride
			return kubeClient.Update(ctx, agent)
		}, "30s", "10s").Should(BeNil())

		Eventually(func() bool {
			condition := conditionsv1.FindStatusCondition(getAgentCRD(ctx, kubeClient, key).Status.Conditions, controllers.SpecSyncedCondition)
			if condition != nil {
				return strings.Contains(condition.Message, "error parsing ignition: config is not valid")
			}
			return false
		}, "15s", "2s").Should(Equal(true))
		h, err = common.GetHostFromDB(db, cluster.ID.String(), host.ID.String())
		Expect(err).To(BeNil())
		Expect(h.IgnitionConfigOverrides).To(BeEmpty())
	})

	It("deploy clusterDeployment with agent and update installer args", func() {
		secretRef := deployLocalObjectSecretIfNeeded(ctx, kubeClient)
		spec := getDefaultClusterDeploymentSpec(secretRef)
		deployClusterDeploymentCRD(ctx, kubeClient, spec)
		key := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      spec.ClusterName,
		}
		cluster := getClusterFromDB(ctx, kubeClient, db, key, waitForReconcileTimeout)
		configureLocalAgentClient(cluster.ID.String())
		host := setupNewHost(ctx, "hostname1", *cluster.ID)
		key = types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      host.ID.String(),
		}

		installerArgs := `["--append-karg", "ip=192.0.2.2::192.0.2.254:255.255.255.0:core0.example.com:enp1s0:none", "--save-partindex", "1", "-n"]`
		Eventually(func() error {
			agent := getAgentCRD(ctx, kubeClient, key)
			agent.Spec.InstallerArgs = installerArgs
			return kubeClient.Update(ctx, agent)
		}, "30s", "10s").Should(BeNil())

		Eventually(func() bool {
			h, err := common.GetHostFromDB(db, cluster.ID.String(), host.ID.String())
			Expect(err).To(BeNil())
			agent := getAgentCRD(ctx, kubeClient, key)

			var j, j2 interface{}
			err = json.Unmarshal([]byte(agent.Spec.InstallerArgs), &j)
			Expect(err).To(BeNil())

			if h.InstallerArgs == "" {
				return false
			}

			err = json.Unmarshal([]byte(h.InstallerArgs), &j2)
			Expect(err).To(BeNil())
			return reflect.DeepEqual(j2, j)
		}, "2m", "10s").Should(Equal(true))

		By("Clean installer args")
		h, err := common.GetHostFromDB(db, cluster.ID.String(), host.ID.String())
		Expect(err).To(BeNil())
		Expect(h.InstallerArgs).NotTo(BeEmpty())

		Eventually(func() error {
			agent := getAgentCRD(ctx, kubeClient, key)
			agent.Spec.InstallerArgs = ""
			return kubeClient.Update(ctx, agent)
		}, "30s", "10s").Should(BeNil())

		Eventually(func() int {
			h, err := common.GetHostFromDB(db, cluster.ID.String(), host.ID.String())
			Expect(err).To(BeNil())
			var j []string
			err = json.Unmarshal([]byte(h.InstallerArgs), &j)
			Expect(err).To(BeNil())

			return len(j)
		}, "2m", "10s").Should(Equal(0))
	})

	It("deploy clusterDeployment with agent and invalid installer args", func() {
		secretRef := deployLocalObjectSecretIfNeeded(ctx, kubeClient)
		spec := getDefaultClusterDeploymentSpec(secretRef)
		deployClusterDeploymentCRD(ctx, kubeClient, spec)
		key := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      spec.ClusterName,
		}
		cluster := getClusterFromDB(ctx, kubeClient, db, key, waitForReconcileTimeout)
		configureLocalAgentClient(cluster.ID.String())
		host := setupNewHost(ctx, "hostname1", *cluster.ID)
		key = types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      host.ID.String(),
		}

		h, err := common.GetHostFromDB(db, cluster.ID.String(), host.ID.String())
		Expect(err).To(BeNil())
		Expect(h.InstallerArgs).To(BeEmpty())

		By("Invalid installer args - invalid json")
		installerArgs := `"--append-karg", "ip=192.0.2.2::192.0.2.254:255.255.255.0:core0.example.com:enp1s0:none", "--save-partindex", "1", "-n"]`
		Eventually(func() error {
			agent := getAgentCRD(ctx, kubeClient, key)
			agent.Spec.InstallerArgs = installerArgs
			return kubeClient.Update(ctx, agent)
		}, "30s", "10s").Should(BeNil())

		Eventually(func() bool {
			condition := conditionsv1.FindStatusCondition(getAgentCRD(ctx, kubeClient, key).Status.Conditions, controllers.SpecSyncedCondition)
			if condition != nil {
				return strings.Contains(condition.Message, "Fail to unmarshal installer args")
			}
			return false
		}, "15s", "2s").Should(Equal(true))
		h, err = common.GetHostFromDB(db, cluster.ID.String(), host.ID.String())
		Expect(err).To(BeNil())
		Expect(h.InstallerArgs).To(BeEmpty())

		By("Invalid installer args - invalid params")
		installerArgs = `["--wrong-param", "ip=192.0.2.2::192.0.2.254:255.255.255.0:core0.example.com:enp1s0:none", "--save-partindex", "1", "-n"]`
		Eventually(func() error {
			agent := getAgentCRD(ctx, kubeClient, key)
			agent.Spec.InstallerArgs = installerArgs
			return kubeClient.Update(ctx, agent)
		}, "30s", "10s").Should(BeNil())

		Eventually(func() bool {
			condition := conditionsv1.FindStatusCondition(getAgentCRD(ctx, kubeClient, key).Status.Conditions, controllers.SpecSyncedCondition)
			if condition != nil {
				return strings.Contains(condition.Message, "found unexpected flag --wrong-param")
			}
			return false
		}, "15s", "2s").Should(Equal(true))
		h, err = common.GetHostFromDB(db, cluster.ID.String(), host.ID.String())
		Expect(err).To(BeNil())
		Expect(h.InstallerArgs).To(BeEmpty())
	})

	It("deploy clusterDeployment with agent,bmh and installer args", func() {
		secretRef := deployLocalObjectSecretIfNeeded(ctx, kubeClient)
		spec := getDefaultClusterDeploymentSpec(secretRef)
		deployClusterDeploymentCRD(ctx, kubeClient, spec)
		key := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      spec.ClusterName,
		}
		cluster := getClusterFromDB(ctx, kubeClient, db, key, waitForReconcileTimeout)
		configureLocalAgentClient(cluster.ID.String())
		host := setupNewHost(ctx, "hostname1", *cluster.ID)
		key = types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      host.ID.String(),
		}

		agent := getAgentCRD(ctx, kubeClient, key)
		bmhSpec := bmhv1alpha1.BareMetalHostSpec{BootMACAddress: getAgentMac(agent)}
		deployBMHCRD(ctx, kubeClient, host.ID.String(), &bmhSpec)

		installerArgs := `["--append-karg", "ip=192.0.2.2::192.0.2.254:255.255.255.0:core0.example.com:enp1s0:none", "--save-partindex", "1", "-n"]`

		Eventually(func() error {
			bmh := getBmhCRD(ctx, kubeClient, key)
			bmh.SetAnnotations(map[string]string{controllers.BMH_AGENT_INSTALLER_ARGS: installerArgs})
			return kubeClient.Update(ctx, bmh)
		}, "30s", "10s").Should(BeNil())

		Eventually(func() bool {
			h, err := common.GetHostFromDB(db, cluster.ID.String(), host.ID.String())
			Expect(err).To(BeNil())
			agent := getAgentCRD(ctx, kubeClient, key)
			if agent.Spec.InstallerArgs == "" {
				return false
			}

			var j, j2 interface{}
			err = json.Unmarshal([]byte(agent.Spec.InstallerArgs), &j)
			Expect(err).To(BeNil())

			if h.InstallerArgs == "" {
				return false
			}

			err = json.Unmarshal([]byte(h.InstallerArgs), &j2)
			Expect(err).To(BeNil())
			return reflect.DeepEqual(j2, j)
		}, "2m", "10s").Should(Equal(true))

		By("Clean installer args")
		h, err := common.GetHostFromDB(db, cluster.ID.String(), host.ID.String())
		Expect(err).To(BeNil())
		Expect(h.InstallerArgs).NotTo(BeEmpty())

		Eventually(func() error {
			bmh := getBmhCRD(ctx, kubeClient, key)
			bmh.SetAnnotations(map[string]string{controllers.BMH_AGENT_INSTALLER_ARGS: ""})
			return kubeClient.Update(ctx, bmh)
		}, "30s", "10s").Should(BeNil())

		By("Verify agent installer args were cleaned")
		Eventually(func() string {
			agent := getAgentCRD(ctx, kubeClient, key)
			return agent.Spec.InstallerArgs
		}, "30s", "10s").Should(BeEmpty())

		By("Verify host installer args were cleaned")
		Eventually(func() int {
			h, err := common.GetHostFromDB(db, cluster.ID.String(), host.ID.String())
			Expect(err).To(BeNil())

			var j []string
			err = json.Unmarshal([]byte(h.InstallerArgs), &j)
			Expect(err).To(BeNil())

			return len(j)
		}, "2m", "10s").Should(Equal(0))
	})

	It("deploy clusterDeployment and infraEnv and verify cluster updates", func() {
		infraEnvName := "infraenv"
		secretRef := deployLocalObjectSecretIfNeeded(ctx, kubeClient)
		clusterDeploymentSpec := getDefaultClusterDeploymentSNOSpec(secretRef)
		deployClusterDeploymentCRD(ctx, kubeClient, clusterDeploymentSpec)
		clusterKubeName := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      clusterDeploymentSpec.ClusterName,
		}
		checkClusterCondition(ctx, clusterKubeName, controllers.ClusterReadyForInstallationCondition, controllers.ClusterNotReadyReason)
		cluster := getClusterFromDB(ctx, kubeClient, db, clusterKubeName, waitForReconcileTimeout)
		configureLocalAgentClient(cluster.ID.String())
		Expect(cluster.NoProxy).Should(Equal(""))
		Expect(cluster.HTTPProxy).Should(Equal(""))
		Expect(cluster.HTTPSProxy).Should(Equal(""))
		Expect(cluster.AdditionalNtpSource).Should(Equal(""))

		infraEnvSpec := getDefaultInfraEnvSpec(secretRef, clusterDeploymentSpec)
		infraEnvSpec.Proxy = &v1beta1.Proxy{
			NoProxy:    "192.168.1.1",
			HTTPProxy:  "http://192.168.1.2",
			HTTPSProxy: "http://192.168.1.3",
		}
		infraEnvSpec.AdditionalNTPSources = []string{"192.168.1.4"}
		deployInfraEnvCRD(ctx, kubeClient, infraEnvName, infraEnvSpec)
		infraEnvKubeName := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      infraEnvName,
		}
		// InfraEnv Reconcile takes longer, since it needs to generate the image.
		checkInfraEnvCondition(ctx, infraEnvKubeName, v1beta1.ImageCreatedCondition, v1beta1.ImageStateCreated)
		cluster = getClusterFromDB(ctx, kubeClient, db, clusterKubeName, waitForReconcileTimeout)
		Expect(cluster.ImageGenerated).Should(Equal(true))
		By("Validate proxy settings.", func() {
			Expect(cluster.NoProxy).Should(Equal("192.168.1.1"))
			Expect(cluster.HTTPProxy).Should(Equal("http://192.168.1.2"))
			Expect(cluster.HTTPSProxy).Should(Equal("http://192.168.1.3"))
		})
		By("Validate additional NTP settings.")
		Expect(cluster.AdditionalNtpSource).Should(ContainSubstring("192.168.1.4"))
		By("InfraEnv image type defaults to minimal-iso.")
		Expect(cluster.ImageInfo.Type).Should(Equal(models.ImageTypeMinimalIso))
	})

	It("deploy clusterDeployment and infraEnv with ignition override", func() {
		infraEnvName := "infraenv"
		secretRef := deployLocalObjectSecretIfNeeded(ctx, kubeClient)
		clusterDeploymentSpec := getDefaultClusterDeploymentSNOSpec(secretRef)
		deployClusterDeploymentCRD(ctx, kubeClient, clusterDeploymentSpec)
		clusterKubeName := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      clusterDeploymentSpec.ClusterName,
		}
		checkClusterCondition(ctx, clusterKubeName, controllers.ClusterReadyForInstallationCondition, controllers.ClusterNotReadyReason)

		cluster := getClusterFromDB(ctx, kubeClient, db, clusterKubeName, waitForReconcileTimeout)
		configureLocalAgentClient(cluster.ID.String())
		Expect(cluster.IgnitionConfigOverrides).Should(Equal(""))

		infraEnvSpec := getDefaultInfraEnvSpec(secretRef, clusterDeploymentSpec)
		infraEnvSpec.IgnitionConfigOverride = fakeIgnitionConfigOverride

		deployInfraEnvCRD(ctx, kubeClient, infraEnvName, infraEnvSpec)
		infraEnvKubeName := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      infraEnvName,
		}
		// InfraEnv Reconcile takes longer, since it needs to generate the image.
		checkInfraEnvCondition(ctx, infraEnvKubeName, v1beta1.ImageCreatedCondition, v1beta1.ImageStateCreated)
		cluster = getClusterFromDB(ctx, kubeClient, db, clusterKubeName, waitForReconcileTimeout)
		Expect(cluster.IgnitionConfigOverrides).Should(Equal(fakeIgnitionConfigOverride))
		Expect(cluster.ImageGenerated).Should(Equal(true))
	})

	It("deploy infraEnv before clusterDeployment", func() {
		infraEnvName := "infraenv"
		secretRef := deployLocalObjectSecretIfNeeded(ctx, kubeClient)
		clusterDeploymentSpec := getDefaultClusterDeploymentSNOSpec(secretRef)
		infraEnvSpec := getDefaultInfraEnvSpec(secretRef, clusterDeploymentSpec)

		deployInfraEnvCRD(ctx, kubeClient, infraEnvName, infraEnvSpec)
		infraEnvKubeName := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      infraEnvName,
		}
		deployClusterDeploymentCRD(ctx, kubeClient, clusterDeploymentSpec)
		clusterKubeName := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      clusterDeploymentSpec.ClusterName,
		}
		checkClusterCondition(ctx, clusterKubeName, controllers.ClusterReadyForInstallationCondition, controllers.ClusterNotReadyReason)
		cluster := getClusterFromDB(ctx, kubeClient, db, clusterKubeName, waitForReconcileTimeout)
		configureLocalAgentClient(cluster.ID.String())

		checkInfraEnvCondition(ctx, infraEnvKubeName, v1beta1.ImageCreatedCondition, v1beta1.ImageStateCreated)
		cluster = getClusterFromDB(ctx, kubeClient, db, clusterKubeName, waitForReconcileTimeout)
		Expect(cluster.ImageGenerated).Should(Equal(true))
	})

	It("deploy clusterDeployment and infraEnv and with an invalid ignition override", func() {
		infraEnvName := "infraenv"
		secretRef := deployLocalObjectSecretIfNeeded(ctx, kubeClient)
		clusterDeploymentSpec := getDefaultClusterDeploymentSNOSpec(secretRef)
		deployClusterDeploymentCRD(ctx, kubeClient, clusterDeploymentSpec)
		clusterKubeName := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      clusterDeploymentSpec.ClusterName,
		}
		checkClusterCondition(ctx, clusterKubeName, controllers.ClusterReadyForInstallationCondition, controllers.ClusterNotReadyReason)
		cluster := getClusterFromDB(ctx, kubeClient, db, clusterKubeName, waitForReconcileTimeout)
		configureLocalAgentClient(cluster.ID.String())
		Expect(cluster.IgnitionConfigOverrides).Should(Equal(""))

		infraEnvSpec := getDefaultInfraEnvSpec(secretRef, clusterDeploymentSpec)
		infraEnvSpec.IgnitionConfigOverride = badIgnitionConfigOverride

		deployInfraEnvCRD(ctx, kubeClient, infraEnvName, infraEnvSpec)
		infraEnvKubeName := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      infraEnvName,
		}
		checkInfraEnvCondition(ctx, infraEnvKubeName, v1beta1.ImageCreatedCondition, v1beta1.ImageStateFailedToCreate+": error parsing ignition: config is not valid")
		cluster = getClusterFromDB(ctx, kubeClient, db, clusterKubeName, waitForReconcileTimeout)
		Expect(cluster.IgnitionConfigOverrides).ShouldNot(Equal(fakeIgnitionConfigOverride))
		Expect(cluster.ImageGenerated).Should(Equal(false))

	})

	It("deploy clusterDeployment with install config override", func() {
		secretRef := deployLocalObjectSecretIfNeeded(ctx, kubeClient)
		clusterDeploymentSpec := getDefaultClusterDeploymentSNOSpec(secretRef)
		deployClusterDeploymentCRD(ctx, kubeClient, clusterDeploymentSpec)
		clusterKubeName := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      clusterDeploymentSpec.ClusterName,
		}
		checkClusterCondition(ctx, clusterKubeName, controllers.ClusterReadyForInstallationCondition, controllers.ClusterNotReadyReason)
		cluster := getClusterFromDB(ctx, kubeClient, db, clusterKubeName, waitForReconcileTimeout)
		Expect(cluster.InstallConfigOverrides).Should(Equal(""))

		installConfigOverrides := `{"controlPlane": {"hyperthreading": "Enabled"}}`
		addAnnotationToClusterDeployment(ctx, kubeClient, clusterKubeName, controllers.InstallConfigOverrides, installConfigOverrides)

		Eventually(func() string {
			c := getClusterFromDB(ctx, kubeClient, db, clusterKubeName, waitForReconcileTimeout)
			if c != nil {
				return c.InstallConfigOverrides
			}
			return ""
		}, "1m", "2s").Should(Equal(installConfigOverrides))
	})

	It("deploy clusterDeployment with malformed install config override", func() {
		secretRef := deployLocalObjectSecretIfNeeded(ctx, kubeClient)
		clusterDeploymentSpec := getDefaultClusterDeploymentSNOSpec(secretRef)
		deployClusterDeploymentCRD(ctx, kubeClient, clusterDeploymentSpec)
		clusterKubeName := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      clusterDeploymentSpec.ClusterName,
		}
		checkClusterCondition(ctx, clusterKubeName, controllers.ClusterReadyForInstallationCondition, controllers.ClusterNotReadyReason)
		cluster := getClusterFromDB(ctx, kubeClient, db, clusterKubeName, waitForReconcileTimeout)
		Expect(cluster.InstallConfigOverrides).Should(Equal(""))

		installConfigOverrides := `{"controlPlane": "malformed json": "Enabled"}}`
		addAnnotationToClusterDeployment(ctx, kubeClient, clusterKubeName, controllers.InstallConfigOverrides, installConfigOverrides)
		checkClusterCondition(ctx, clusterKubeName, controllers.ClusterSpecSyncedCondition, controllers.InputErrorReason)
		cluster = getClusterFromDB(ctx, kubeClient, db, clusterKubeName, waitForReconcileTimeout)
		Expect(cluster.InstallConfigOverrides).Should(Equal(""))
	})

	It("deploy clusterDeployment and infraEnv and with NMState config", func() {
		var (
			NMStateLabelName  = "someName"
			NMStateLabelValue = "someValue"
			nicPrimary        = "eth0"
			nicSecondary      = "eth1"
			macPrimary        = "09:23:0f:d8:92:AA"
			macSecondary      = "09:23:0f:d8:92:AB"
			ip4Primary        = "192.168.126.30"
			ip4Secondary      = "192.168.140.30"
			dnsGW             = "192.168.126.1"
		)
		hostStaticNetworkConfig := common.FormatStaticConfigHostYAML(
			nicPrimary, nicSecondary, ip4Primary, ip4Secondary, dnsGW,
			models.MacInterfaceMap{
				{MacAddress: macPrimary, LogicalNicName: nicPrimary},
				{MacAddress: macSecondary, LogicalNicName: nicSecondary},
			})
		nmstateConfigSpec := getDefaultNMStateConfigSpec(nicPrimary, nicSecondary, macPrimary, macSecondary, hostStaticNetworkConfig.NetworkYaml)
		deployNMStateConfigCRD(ctx, kubeClient, "nmstate1", NMStateLabelName, NMStateLabelValue, nmstateConfigSpec)
		infraEnvName := "infraenv"
		secretRef := deployLocalObjectSecretIfNeeded(ctx, kubeClient)
		clusterDeploymentSpec := getDefaultClusterDeploymentSNOSpec(secretRef)
		deployClusterDeploymentCRD(ctx, kubeClient, clusterDeploymentSpec)
		clusterKubeName := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      clusterDeploymentSpec.ClusterName,
		}
		checkClusterCondition(ctx, clusterKubeName, controllers.ClusterReadyForInstallationCondition, controllers.ClusterNotReadyReason)
		infraEnvSpec := getDefaultInfraEnvSpec(secretRef, clusterDeploymentSpec)
		infraEnvSpec.NMStateConfigLabelSelector = metav1.LabelSelector{MatchLabels: map[string]string{NMStateLabelName: NMStateLabelValue}}
		deployInfraEnvCRD(ctx, kubeClient, infraEnvName, infraEnvSpec)
		infraEnvKubeName := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      infraEnvName,
		}
		// InfraEnv Reconcile takes longer, since it needs to generate the image.
		checkInfraEnvCondition(ctx, infraEnvKubeName, v1beta1.ImageCreatedCondition, v1beta1.ImageStateCreated)
		cluster := getClusterFromDB(ctx, kubeClient, db, clusterKubeName, waitForReconcileTimeout)
		Expect(cluster.ImageInfo.StaticNetworkConfig).Should(ContainSubstring(hostStaticNetworkConfig.NetworkYaml))
		Expect(cluster.ImageGenerated).Should(Equal(true))
	})

	It("deploy clusterDeployment and infraEnv and with an invalid NMState config YAML", func() {
		Skip("MGMT-5324 flaky test")
		var (
			NMStateLabelName  = "someName"
			NMStateLabelValue = "someValue"
			nicPrimary        = "eth0"
			nicSecondary      = "eth1"
			macPrimary        = "09:23:0f:d8:92:AA"
			macSecondary      = "09:23:0f:d8:92:AB"
		)
		nmstateConfigSpec := getDefaultNMStateConfigSpec(nicPrimary, nicSecondary, macPrimary, macSecondary, "foo: bar")
		deployNMStateConfigCRD(ctx, kubeClient, "nmstate2", NMStateLabelName, NMStateLabelValue, nmstateConfigSpec)
		infraEnvName := "infraenv"
		secretRef := deployLocalObjectSecretIfNeeded(ctx, kubeClient)
		clusterDeploymentSpec := getDefaultClusterDeploymentSNOSpec(secretRef)
		deployClusterDeploymentCRD(ctx, kubeClient, clusterDeploymentSpec)
		clusterKubeName := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      clusterDeploymentSpec.ClusterName,
		}
		checkClusterCondition(ctx, clusterKubeName, controllers.ClusterReadyForInstallationCondition, controllers.ClusterNotReadyReason)
		infraEnvSpec := getDefaultInfraEnvSpec(secretRef, clusterDeploymentSpec)
		infraEnvSpec.NMStateConfigLabelSelector = metav1.LabelSelector{MatchLabels: map[string]string{NMStateLabelName: NMStateLabelValue}}
		deployInfraEnvCRD(ctx, kubeClient, infraEnvName, infraEnvSpec)
		infraEnvKubeName := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      infraEnvName,
		}
		// InfraEnv Reconcile takes longer, since it needs to generate the image.
		checkInfraEnvCondition(ctx, infraEnvKubeName, v1beta1.ImageCreatedCondition, fmt.Sprintf("%s: internal error", v1beta1.ImageStateFailedToCreate))
		cluster := getClusterFromDB(ctx, kubeClient, db, clusterKubeName, waitForReconcileTimeout)
		Expect(cluster.ImageGenerated).Should(Equal(false))
	})

	It("SNO deploy clusterDeployment full install and validate MetaData", func() {
		By("Create cluster")
		secretRef := deployLocalObjectSecretIfNeeded(ctx, kubeClient)
		spec := getDefaultClusterDeploymentSNOSpec(secretRef)
		deployClusterDeploymentCRD(ctx, kubeClient, spec)
		clusterKey := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      spec.ClusterName,
		}
		cluster := getClusterFromDB(ctx, kubeClient, db, clusterKey, waitForReconcileTimeout)
		configureLocalAgentClient(cluster.ID.String())
		host := setupNewHost(ctx, "hostname1", *cluster.ID)
		key := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      host.ID.String(),
		}
		By("Approve Agent")
		Eventually(func() error {
			agent := getAgentCRD(ctx, kubeClient, key)
			agent.Spec.Approved = true
			return kubeClient.Update(ctx, agent)
		}, "30s", "10s").Should(BeNil())

		By("Wait for installing")
		checkClusterCondition(ctx, clusterKey, controllers.ClusterInstalledCondition, controllers.InstallationInProgressReason)

		Eventually(func() bool {
			c := getClusterFromDB(ctx, kubeClient, db, clusterKey, waitForReconcileTimeout)
			for _, h := range c.Hosts {
				if !funk.ContainsString([]string{models.HostStatusInstalling, models.HostStatusDisabled}, swag.StringValue(h.Status)) {
					return false
				}
			}
			return true
		}, "1m", "2s").Should(BeTrue())

		updateProgress(*host.ID, *cluster.ID, models.HostStageDone)

		By("Complete Installation")
		completeInstallation(agentBMClient, *cluster.ID)
		isSuccess := true
		_, err := agentBMClient.Installer.CompleteInstallation(ctx, &installer.CompleteInstallationParams{
			ClusterID: *cluster.ID,
			CompletionParams: &models.CompletionParams{
				IsSuccess: &isSuccess,
			},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Verify Cluster Metadata")
		Eventually(func() bool {
			return getClusterDeploymentCRD(ctx, kubeClient, clusterKey).Spec.Installed
		}, "1m", "2s").Should(BeTrue())
		Expect(getClusterDeploymentCRD(ctx, kubeClient, clusterKey).Status.APIURL).Should(Equal("https://api.test-cluster-sno.hive.example.com:6443"))
		Expect(getClusterDeploymentCRD(ctx, kubeClient, clusterKey).Status.WebConsoleURL).Should(Equal("https://console-openshift-console.apps.test-cluster-sno.hive.example.com"))
		passwordSecretRef := getClusterDeploymentCRD(ctx, kubeClient, clusterKey).Spec.ClusterMetadata.AdminPasswordSecretRef
		Expect(passwordSecretRef).NotTo(BeNil())
		passwordkey := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      passwordSecretRef.Name,
		}
		passwordSecret := getSecret(ctx, kubeClient, passwordkey)
		Expect(passwordSecret.Data["password"]).NotTo(BeNil())
		Expect(passwordSecret.Data["username"]).NotTo(BeNil())
		configSecretRef := getClusterDeploymentCRD(ctx, kubeClient, clusterKey).Spec.ClusterMetadata.AdminKubeconfigSecretRef
		Expect(passwordSecretRef).NotTo(BeNil())
		configkey := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      configSecretRef.Name,
		}
		configSecret := getSecret(ctx, kubeClient, configkey)
		Expect(configSecret.Data["kubeconfig"]).NotTo(BeNil())
	})

	It("None SNO deploy clusterDeployment full install and validate MetaData", func() {
		By("Create cluster")
		secretRef := deployLocalObjectSecretIfNeeded(ctx, kubeClient)
		spec := getDefaultClusterDeploymentSpec(secretRef)
		deployClusterDeploymentCRD(ctx, kubeClient, spec)
		clusterKey := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      spec.ClusterName,
		}
		cluster := getClusterFromDB(ctx, kubeClient, db, clusterKey, waitForReconcileTimeout)
		configureLocalAgentClient(cluster.ID.String())
		hosts := make([]*models.Host, 0)
		for i := 0; i < 3; i++ {
			hostname := fmt.Sprintf("h%d", i)
			host := setupNewHost(ctx, hostname, *cluster.ID)
			hosts = append(hosts, host)
		}
		for _, host := range hosts {
			checkAgentCondition(ctx, host.ID.String(), controllers.ValidatedCondition, controllers.ValidationsFailingReason)
		}
		generateFullMeshConnectivity(ctx, "1.2.3.10", hosts...)
		for _, host := range hosts {
			checkAgentCondition(ctx, host.ID.String(), controllers.ValidatedCondition, controllers.ValidationsPassingReason)
		}

		By("Approve Agents")
		for _, host := range hosts {
			hostkey := types.NamespacedName{
				Namespace: Options.Namespace,
				Name:      host.ID.String(),
			}
			Eventually(func() error {
				agent := getAgentCRD(ctx, kubeClient, hostkey)
				agent.Spec.Approved = true
				return kubeClient.Update(ctx, agent)
			}, "30s", "10s").Should(BeNil())
		}

		By("Wait for installing")
		checkClusterCondition(ctx, clusterKey, controllers.ClusterInstalledCondition, controllers.InstallationInProgressReason)
		Eventually(func() bool {
			c := getClusterFromDB(ctx, kubeClient, db, clusterKey, waitForReconcileTimeout)
			for _, h := range c.Hosts {
				if !funk.ContainsString([]string{models.HostStatusInstalling, models.HostStatusDisabled}, swag.StringValue(h.Status)) {
					return false
				}
			}
			return true
		}, "1m", "2s").Should(BeTrue())

		for _, host := range hosts {
			checkAgentCondition(ctx, host.ID.String(), controllers.InstalledCondition, controllers.InstallationInProgressReason)
		}

		for _, host := range hosts {
			updateProgress(*host.ID, *cluster.ID, models.HostStageDone)
		}

		By("Complete Installation")
		completeInstallation(agentBMClient, *cluster.ID)
		isSuccess := true
		_, err := agentBMClient.Installer.CompleteInstallation(ctx, &installer.CompleteInstallationParams{
			ClusterID: *cluster.ID,
			CompletionParams: &models.CompletionParams{
				IsSuccess: &isSuccess,
			},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Verify Day 2 Cluster")
		checkClusterCondition(ctx, clusterKey, controllers.ClusterInstalledCondition, controllers.InstalledReason)
		cluster = getClusterFromDB(ctx, kubeClient, db, clusterKey, waitForReconcileTimeout)
		Expect(*cluster.Kind).Should(Equal(models.ClusterKindAddHostsCluster))

		By("Verify Cluster Metadata")
		Eventually(func() bool {
			return getClusterDeploymentCRD(ctx, kubeClient, clusterKey).Spec.Installed
		}, "2m", "2s").Should(BeTrue())
		passwordSecretRef := getClusterDeploymentCRD(ctx, kubeClient, clusterKey).Spec.ClusterMetadata.AdminPasswordSecretRef
		Expect(passwordSecretRef).NotTo(BeNil())
		passwordkey := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      passwordSecretRef.Name,
		}
		passwordSecret := getSecret(ctx, kubeClient, passwordkey)
		Expect(passwordSecret.Data["password"]).NotTo(BeNil())
		Expect(passwordSecret.Data["username"]).NotTo(BeNil())
		configSecretRef := getClusterDeploymentCRD(ctx, kubeClient, clusterKey).Spec.ClusterMetadata.AdminKubeconfigSecretRef
		Expect(passwordSecretRef).NotTo(BeNil())
		configkey := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      configSecretRef.Name,
		}
		configSecret := getSecret(ctx, kubeClient, configkey)
		Expect(configSecret.Data["kubeconfig"]).NotTo(BeNil())
	})

	It("None SNO deploy clusterDeployment full install and Day 2 new host", func() {
		By("Create cluster")
		secretRef := deployLocalObjectSecretIfNeeded(ctx, kubeClient)
		spec := getDefaultClusterDeploymentSpec(secretRef)
		deployClusterDeploymentCRD(ctx, kubeClient, spec)
		clusterKey := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      spec.ClusterName,
		}
		cluster := getClusterFromDB(ctx, kubeClient, db, clusterKey, waitForReconcileTimeout)
		configureLocalAgentClient(cluster.ID.String())
		hosts := make([]*models.Host, 0)
		for i := 0; i < 3; i++ {
			hostname := fmt.Sprintf("h%d", i)
			host := setupNewHost(ctx, hostname, *cluster.ID)
			hosts = append(hosts, host)
		}
		generateFullMeshConnectivity(ctx, "1.2.3.10", hosts...)
		By("Approve Agents")
		for _, host := range hosts {
			hostkey := types.NamespacedName{
				Namespace: Options.Namespace,
				Name:      host.ID.String(),
			}
			Eventually(func() error {
				agent := getAgentCRD(ctx, kubeClient, hostkey)
				agent.Spec.Approved = true
				return kubeClient.Update(ctx, agent)
			}, "30s", "10s").Should(BeNil())
		}

		By("Wait for installing")
		checkClusterCondition(ctx, clusterKey, controllers.ClusterInstalledCondition, controllers.InstallationInProgressReason)
		Eventually(func() bool {
			c := getClusterFromDB(ctx, kubeClient, db, clusterKey, waitForReconcileTimeout)
			for _, h := range c.Hosts {
				if !funk.ContainsString([]string{models.HostStatusInstalling, models.HostStatusDisabled}, swag.StringValue(h.Status)) {
					return false
				}
			}
			return true
		}, "1m", "2s").Should(BeTrue())

		for _, host := range hosts {
			updateProgress(*host.ID, *cluster.ID, models.HostStageDone)
		}

		By("Complete Installation")
		completeInstallation(agentBMClient, *cluster.ID)
		isSuccess := true
		_, err := agentBMClient.Installer.CompleteInstallation(ctx, &installer.CompleteInstallationParams{
			ClusterID: *cluster.ID,
			CompletionParams: &models.CompletionParams{
				IsSuccess: &isSuccess,
			},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Verify Day 2 Cluster")
		checkClusterCondition(ctx, clusterKey, controllers.ClusterInstalledCondition, controllers.InstalledReason)
		cluster = getClusterFromDB(ctx, kubeClient, db, clusterKey, waitForReconcileTimeout)
		Expect(*cluster.Kind).Should(Equal(models.ClusterKindAddHostsCluster))

		By("Add Day 2 host and approve agent")
		configureLocalAgentClient(cluster.ID.String())
		host := setupNewHost(ctx, "hostnameday2", *cluster.ID)
		key := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      host.ID.String(),
		}
		generateApiVipPostStepReply(ctx, host, true)
		Eventually(func() error {
			agent := getAgentCRD(ctx, kubeClient, key)
			agent.Spec.Approved = true
			return kubeClient.Update(ctx, agent)
		}, "30s", "10s").Should(BeNil())

		checkAgentCondition(ctx, host.ID.String(), controllers.InstalledCondition, controllers.InstallationInProgressReason)
		checkAgentCondition(ctx, host.ID.String(), controllers.ReadyForInstallationCondition, controllers.AgentAlreadyInstallingReason)
		checkAgentCondition(ctx, host.ID.String(), controllers.SpecSyncedCondition, controllers.SyncedOkReason)
		checkAgentCondition(ctx, host.ID.String(), controllers.ConnectedCondition, controllers.AgentConnectedReason)
		checkAgentCondition(ctx, host.ID.String(), controllers.ValidatedCondition, controllers.ValidationsPassingReason)
	})

	It("deploy clusterDeployment with invalid machine cidr", func() {
		secretRef := deployLocalObjectSecretIfNeeded(ctx, kubeClient)
		clusterDeploymentSpec := getDefaultClusterDeploymentSNOSpec(secretRef)
		clusterDeploymentSpec.Provisioning.InstallStrategy.Agent.Networking.MachineNetwork = []agentv1.MachineNetworkEntry{{CIDR: "1.2.3.5/24"}}
		deployClusterDeploymentCRD(ctx, kubeClient, clusterDeploymentSpec)
		clusterKubeName := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      clusterDeploymentSpec.ClusterName,
		}
		checkClusterCondition(ctx, clusterKubeName, controllers.ClusterReadyForInstallationCondition, controllers.ClusterNotReadyReason)
	})

	It("deploy clusterDeployment without machine cidr", func() {
		secretRef := deployLocalObjectSecretIfNeeded(ctx, kubeClient)
		clusterDeploymentSpec := getDefaultClusterDeploymentSNOSpec(secretRef)
		clusterDeploymentSpec.Provisioning.InstallStrategy.Agent.Networking.MachineNetwork = []agentv1.MachineNetworkEntry{}
		deployClusterDeploymentCRD(ctx, kubeClient, clusterDeploymentSpec)
		clusterKubeName := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      clusterDeploymentSpec.ClusterName,
		}
		checkClusterCondition(ctx, clusterKubeName, controllers.ClusterReadyForInstallationCondition, controllers.ClusterNotReadyReason)
	})

	It("deploy clusterDeployment with invalid clusterImageSet", func() {
		secretRef := deployLocalObjectSecretIfNeeded(ctx, kubeClient)
		clusterDeploymentSpec := getDefaultClusterDeploymentSpec(secretRef)
		clusterDeploymentSpec.Provisioning.ImageSetRef.Name = "invalid"
		deployClusterDeploymentCRD(ctx, kubeClient, clusterDeploymentSpec)
		clusterKubeName := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      clusterDeploymentSpec.ClusterName,
		}
		checkClusterCondition(ctx, clusterKubeName, controllers.ClusterSpecSyncedCondition, controllers.BackendErrorReason)
	})

	It("deploy clusterDeployment with missing clusterImageSet", func() {
		secretRef := deployLocalObjectSecretIfNeeded(ctx, kubeClient)
		spec := getDefaultClusterDeploymentSpec(secretRef)

		// Create ClusterDeployment
		err := kubeClient.Create(ctx, &hivev1.ClusterDeployment{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ClusterDeployment",
				APIVersion: getAPIVersion(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: Options.Namespace,
				Name:      spec.ClusterName,
			},
			Spec: *spec,
		})
		Expect(err).To(BeNil())

		clusterKubeName := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      spec.ClusterName,
		}
		checkClusterCondition(ctx, clusterKubeName, controllers.ClusterSpecSyncedCondition, controllers.BackendErrorReason)

		// Create ClusterImageSet
		deployClusterImageSetCRD(ctx, kubeClient, spec.Provisioning.ImageSetRef)
		checkClusterCondition(ctx, clusterKubeName, controllers.ClusterSpecSyncedCondition, controllers.SyncedOkReason)
	}) */

	It("deploy clusterDeployment with agent role worker and set spoke BMH", func() {
		secretRef := deployLocalObjectSecretIfNeeded(ctx, kubeClient)
		spec := getDefaultClusterDeploymentSpec(secretRef)
		// While reconciling SpokeBMH, if the cluster is not installed yet, we can't get kubeconfig for the cluster yet.
		spec.Installed = true
		deployClusterDeploymentCRD(ctx, kubeClient, spec)
		key := types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      spec.ClusterName,
		}
		cluster := getClusterFromDB(ctx, kubeClient, db, key, waitForReconcileTimeout)
		configureLocalAgentClient(cluster.ID.String())
		host := setupNewHost(ctx, "hostname1", *cluster.ID)
		key = types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      host.ID.String(),
		}

		agent := getAgentCRD(ctx, kubeClient, key)
		Eventually(func() error {
			// Only worker role is supported for day2 operation
			agent.Spec.Role = "worker"
			return kubeClient.Update(ctx, agent)
		}, "30s", "10s").Should(BeNil())

		bmhSpec := bmhv1alpha1.BareMetalHostSpec{BootMACAddress: getAgentMac(agent)}
		fmt.Printf("\n************ Deploying BMH with ID/Name %s \n", host.ID.String())
		deployBMHCRD(ctx, kubeClient, host.ID.String(), &bmhSpec)

		// Secret contains kubeconfig for the spoke cluster
		secretName := fmt.Sprintf(adminKubeConfigStringTemplate, spec.ClusterName)
		fmt.Printf("\n***************** Trying to create a secret %s/%s", Options.Namespace, secretName)
		Eventually(func() error {
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      secretName,
					Namespace: Options.Namespace,
				},
				Data: map[string][]byte{
					"kubeconfig": []byte(BASIC_KUBECONFIG),
				},
			}
			return kubeClient.Create(ctx, secret)
		}, "30s", "10s").Should(BeNil())

		key = types.NamespacedName{
			Namespace: Options.Namespace,
			Name:      secretName,
		}
		configSecret := getSecret(ctx, kubeClient, key)

		kubeconfigData := configSecret.Data["kubeconfig"]
		// kubeconfigData := []byte(BASIC_KUBECONFIG)
		fmt.Println("*******kubeconfigData**********")
		fmt.Println(string(kubeconfigData))
		Expect(kubeconfigData).NotTo(BeNil())

		clientConfig, err := clientcmd.NewClientConfigFromBytes(kubeconfigData)
		// clientConfig, err := clientcmd.NewClientConfigFromBytes([]byte(BASIC_KUBECONFIG))
		fmt.Println("*******clientConfig**********")
		fmt.Println(clientConfig)
		Expect(err).To(BeNil())

		restConfig, err := clientConfig.ClientConfig()
		fmt.Println("*******restConfig**********")
		fmt.Println(restConfig)
		Expect(restConfig).NotTo(BeNil())
		Expect(err).To(BeNil())

		spokeClient, err := k8sclient.New(restConfig, k8sclient.Options{Scheme: scheme.Scheme})
		fmt.Println("*******spokeClient**********")
		fmt.Println(spokeClient)
		Expect(spokeClient).NotTo(BeNil())
		Expect(err).To(BeNil())

		By("ensureSpokeMachine is created")
		Eventually(func() error {
			key = types.NamespacedName{
				Namespace: Options.Namespace,
				Name:      host.ID.String(),
			}
			bmh := getBmhCRD(ctx, kubeClient, key)
			// machineName := fmt.Sprintf("%s-%s", spec.ClusterName, bmh.Name)
			machineName := bmh.Name
			fmt.Printf("machineName = %s", machineName)

			spokeBMH := &bmhv1alpha1.BareMetalHost{}
			// spokeClient := bmhr.spokeClient
			// err = spokeClient.Get(ctx, types.NamespacedName{Name: host.Name, Namespace: testNamespace}, spokeBMH)

			// machine := &machinev1beta1.Machine{
			// 	ObjectMeta: metav1.ObjectMeta{
			// 		Name:      machineName,
			// 		Namespace: bmh.Namespace,
			// 	},
			// }
			spokeClient.Get(ctx, types.NamespacedName{Name: machineName, Namespace: Options.Namespace}, spokeBMH)
			fmt.Println("\n*** spokeBMH is")
			fmt.Println(spokeBMH)
			return spokeClient.Get(ctx, types.NamespacedName{Name: machineName, Namespace: Options.Namespace}, spokeBMH)
			// return spokeClient.Get(ctx, types.NamespacedName{Name: machineName, Namespace: bmh.Namespace}, spokeBMH)
		}, "30s", "10s").ShouldNot(BeNil())

		By("Clean spoke cluster secret")
		Eventually(func() error {
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      secretName,
					Namespace: Options.Namespace,
				},
			}
			return kubeClient.Delete(ctx, secret)
		}, "30s", "10s").Should(BeNil())
	})
})
