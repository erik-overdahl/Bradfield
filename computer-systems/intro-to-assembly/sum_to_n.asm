section .text
global sum_to_n

sum_to_n:
	mov		rbx, rdi
	mov		rax, rdi
	inc		rax					; rax = n+1
	mul		rbx					; rax = n * (n+1)
	shr		rax, 1				; rax = n*(n+1)/2
	ret
